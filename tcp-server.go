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
	"log"
	"bytes"
)

var count int = 0
var allocs int = 0
var errors int = 0
var cached int = 0
var accepts int = 0
var releases int = 0
var total int64 = 0

func initLog() {

	logfile := "1.log"

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}

	w := io.MultiWriter(os.Stderr, file)
	//	w := io.MultiWriter(file)
	log.SetOutput(w)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	fmt.Println("logging started to ", logfile)
}

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
		p.b = make([]byte, 4*4096)
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
	bad_cnt := 0
	last_cnt := 0
	cycles := 0
	
	for {
		_, err := p.r.Read(p.b)
		if err != nil {
			if err != io.EOF {
				log.Println(num, err)
				errors++
			}
			break
		}
		
		if bytes.HasPrefix(p.b, Capprovider) {
			fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(p.w, "Content-Type: text/plain\r\n")
			fmt.Fprintf(p.w, "Connection: keep-alive\r\n")
			fmt.Fprintf(p.w, "Keep-Alive: timeout=120, max=50000\r\n")
			fmt.Fprintf(p.w, "Content-Length: 1\r\n\r\n\n")
//			fmt.Println("num:", num, cnt)
			cnt++
			last_cnt++
		} else if bytes.HasPrefix(p.b, Cready) {
			fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(p.w, "Content-Type: text/plain\r\n")
			fmt.Fprintf(p.w, "Connection: keep-alive\r\n")
			fmt.Fprintf(p.w, "Keep-Alive: timeout=120, max=50000\r\n")
			fmt.Fprintf(p.w, "Content-Length: 2\r\n\r\n1\n")
			ready_cnt++
		} else {
			log.Println(num, "bad url: ", string(p.b))
			fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(p.w, "Content-Type: text/plain\r\n")
			fmt.Fprintf(p.w, "Connection: keep-alive\r\n")
			fmt.Fprintf(p.w, "Keep-Alive: timeout=120, max=50000\r\n")
			fmt.Fprintf(p.w, "Content-Length: 1\r\n\r\n\n")
			bad_cnt++
		}
		
		p.w.Flush()
		p.r.Reset(conn)
		p.w.Reset(conn)
		time.Sleep(time.Millisecond)

		cycles++
		if cycles%2000 == 0 {
			log.Println(num, "2k cycle. cnt", cnt, "r", ready_cnt, "b", bad_cnt, "c", cycles)
			total += int64(last_cnt)
			last_cnt = 0
		}
	}
	conn.Close()
	releases++
	log.Println(num, "close: ", accepts-releases, "errors", errors)
	total += int64(last_cnt)
}

var m1 runtime.MemStats
var m2 runtime.MemStats
var prev_total int64 = 0

func dumpMemStat() {
	runtime.ReadMemStats(&m2)
	log.Printf("%s M %d %d (%d), F %d %d (%d), a: %d, c: %d, gon: %d, total: %d (%d)\n",
		time.Now(),
		m1.Mallocs, m2.Mallocs, m2.Mallocs - m1.Mallocs,
		m1.Frees, m2.Frees, m2.Frees - m1.Frees,
		allocs, cached, runtime.NumGoroutine(), total, (total - prev_total) / 5)
	runtime.ReadMemStats(&m1)
	prev_total = total
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

var rate_per_sec = 1
var throttle = time.Tick(time.Duration(1e9 / rate_per_sec))


func main() {
	initLog()

	ticker := time.NewTicker(time.Second*5)
	go memStat(ticker)

	ln, err := net.Listen("tcp", ":81")
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        
	log.Println("started")
	num := 0

        for {
                conn, err := ln.Accept()
                accepts++
		log.Println("accepted", accepts, err)

                if err != nil {
			log.Println(err.Error())
                        return
                }
                go handleConnection(conn, num)
                num++
        }
}
