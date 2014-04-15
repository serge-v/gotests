package main

import (
	"fmt"
	"net"
	"runtime"
//	"runtime/pprof"
	"bufio"
	"time"
	"io"
	"os"
	"bytes"
)

var count int = 0
var allocs int = 0
var cached int = 0

type holder struct {
	r *bufio.Reader
	w *bufio.Writer
	b []byte
}

var holderChan = make(chan *holder, 4096)

func getH(conn net.Conn) *holder {
	select {
	case p := <- holderChan:
		p.r.Reset(conn)
		p.w.Reset(conn)
		cached++
		return p
	default:
		p := new(holder)
		p.r = bufio.NewReader(conn)
		p.w = bufio.NewWriterSize(conn, 4096)
		p.b = make([]byte, 4096)
		allocs++
		return p
	}
}

func putH(p *holder) {
	p.r.Reset(nil)
	p.w.Reset(nil)
	select {
	case holderChan <- p:
	default:
	}
}

var Capprovider = []byte("GET /data_provider/appnexus?")
var Cready = []byte("GET /ready.ashx")
var Cip = []byte("ip=")

func handleConnection(conn net.Conn, num int) {
	p := getH(conn)
	defer putH(p)
	
	cnt := 0
	ready_cnt := 0
	
	for {
		n, err := p.r.Read(p.b)
		if err != nil {
			if err != io.EOF {
				fmt.Println(n, err)
			}
			break
		}
		
		if bytes.HasPrefix(p.b, Capprovider) {
			fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n")
	//		fmt.Fprintf(p.w, "Content-Type: text/plain\r\n")
	//		fmt.Fprintf(p.w, "Connection: keep-alive\r\n")
	//		fmt.Fprintf(p.w, "Keep-Alive: timeout=120, max=50000\r\n")
			q_pos := len(Capprovider)
			ip_pos := bytes.Index(p.b[q_pos:], Cip)
			if ip_pos >= 0 {
				amp_pos := bytes.IndexAny(p.b[q_pos+ip_pos:], "& ")
				fmt.Fprintf(p.w, "Content-Length: %d\r\n", amp_pos+2)
				fmt.Fprintf(p.w, "\r\n")
				fmt.Fprintf(p.w, "\nip_%s\n", p.b[q_pos+ip_pos+3:q_pos+ip_pos+amp_pos])
/*				fmt.Println("q_pos", q_pos)
				fmt.Println("ip_pos", ip_pos)
				fmt.Println("amp_pos", amp_pos) */
			} else {
				fmt.Println(n, err, string(p.b))
				os.Exit(1)
			}
			cnt++
		} else if bytes.HasPrefix(p.b, Cready) {
			fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n\r\n1")
			ready_cnt++
		} else {
			fmt.Println(n, err, string(p.b))
			os.Exit(1)
		}
		
		p.w.Flush()
		p.r.Reset(conn)
		p.w.Reset(conn)
		
		if cnt%1000 == 0 {
			fmt.Println("-", num, cnt, ready_cnt)
		}
	}
	conn.Close()
	conns--
	if conns%500 == 0 {
		fmt.Println("close", conns)
	}
}

var conns int = 0
var m1 runtime.MemStats
var m2 runtime.MemStats

func dumpMemStat() {
	runtime.ReadMemStats(&m2)
	fmt.Printf("M %d %d (%d), F %d %d (%d), a: %d, c: %d\n",
		m1.Mallocs, m2.Mallocs, m2.Mallocs - m1.Mallocs,
		m1.Frees, m2.Frees, m2.Frees - m1.Frees,
		allocs, cached)
	runtime.ReadMemStats(&m1)
}
/*
func dumpHeapProfile() {
	p := pprof.Lookup("heap")
	if err := p.WriteTo(os.Stdout, 1); err != nil {
		fmt.Println("heap:", err.Error())
	}
}
*/
func memStat(ticker *time.Ticker) {
	for _ = range(ticker.C) {
//		dumpHeapProfile()
		dumpMemStat()
	}
}

var rate_per_sec = 0.1
var throttle = time.Tick(time.Duration(1e9 / rate_per_sec))


func main() {
	
	ticker := time.NewTicker(time.Second*5)
	go memStat(ticker)
	
	ln, err := net.Listen("tcp", ":8080")
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        
	fmt.Println("started")
	num := 0

        for {
                conn, err := ln.Accept()
                conns++
                fmt.Println("accepted", conns, err)
                if err != nil {
			fmt.Println(err.Error())
                        return
                }
                go handleConnection(conn, num)
                num++
                if conns > 10 {
			<-throttle
		}
        }
}
