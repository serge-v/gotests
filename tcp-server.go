package main

import (
	"fmt"
	"net"
	"runtime"
	//	"runtime/pprof"
	"./util"
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

// reusable buffer
type rbuf struct {
	r *bufio.Reader
	w *bufio.Writer
	b []byte
}

type Server struct {
	port             int
	rbufChan         chan *rbuf // reusable buffers cache
	dataEndpoint     []byte
	readyEndpoint    []byte
	stopFile         string
	stopping         bool
	keepaliveMax     int
	keepaliveTimeout int

	// statistics
	connAllocs   int64 // allocated buffers
	connCached   int64 // get from cache count
	connReleases int64 // released count
	tcpAccepts   int64 // TODO: concurrent
	errors       int64 // TODO: concurrent
	totalData    int64 // TODO: concurrent
	totalReady   int64 // TODO: concurrent
}

func newServer(port int) *Server {
	s := Server{
		port:             port,
		stopFile:         fmt.Sprintf("stop-%d.txt", port),
		rbufChan:         make(chan *rbuf, 4096),
		dataEndpoint:     []byte("GET /data_provider/appnexus?"),
		readyEndpoint:    []byte("GET /ready.ashx"),
		keepaliveMax:     20000,
		keepaliveTimeout: 1200,
	}
	return &s
}

func (s *Server) claimBuffer(conn net.Conn) *rbuf {
	select {
	case p := <-s.rbufChan:
		p.r.Reset(conn)
		p.w.Reset(conn)
		s.connCached++
		return p
	default:
		p := new(rbuf)
		p.r = bufio.NewReader(conn)
		p.w = bufio.NewWriterSize(conn, 4096)
		p.b = make([]byte, 4*4096)
		s.connAllocs++
		return p
	}
}

func (s *Server) releaseBuffer(p *rbuf) {
	p.r.Reset(nil)
	p.w.Reset(nil)
	select {
	case s.rbufChan <- p:
		s.connReleases++
	default:
	}
}

func (s *Server) checkStoppingFlag() {
	_, err := os.Stat(s.stopFile)
	s.stopping = (err == nil)
}

func (s *Server) handleConnection(conn net.Conn, num int) {
	p := s.claimBuffer(conn)
	defer s.releaseBuffer(p)

	cnt := 0
	ready_cnt := 0
	bad_cnt := 0
	last_cnt := 0
	cycles := 0
	start := time.Now()
	timeout := time.Duration(time.Second * time.Duration(s.keepaliveTimeout+num%120)) // distribute along 2 minutes
	send_close := false

	for {
		_, err := p.r.Read(p.b)
		if err != nil {
			if err != io.EOF {
				s.errors++
			}
			break
		}

		if (cycles >= s.keepaliveMax-1) || (time.Since(start) > timeout) {
			send_close = true
		}

		fmt.Fprintf(p.w, "HTTP/1.1 200 OK\r\n")
		fmt.Fprintf(p.w, "Content-Type: text/plain\r\n")

		if send_close {
			fmt.Fprintf(p.w, "Connection: close\r\n")
		}

		if bytes.HasPrefix(p.b, s.dataEndpoint) {
			fmt.Fprintf(p.w, "Content-Length: 1\r\n\r\n\n")
			cnt++
			last_cnt++
		} else if bytes.HasPrefix(p.b, s.readyEndpoint) {
			if s.stopping {
				fmt.Fprintf(p.w, "Content-Length: 2\r\n\r\n0\n")
			} else {
				fmt.Fprintf(p.w, "Content-Length: 2\r\n\r\n1\n")
			}
			ready_cnt++
			s.totalReady++
		} else {
			//			log.Println(num, "bad url: ", string(p.b))
			fmt.Fprintf(p.w, "Content-Length: 1\r\n\r\n\n")
			bad_cnt++
		}

		p.w.Flush()
		p.r.Reset(conn)
		p.w.Reset(conn)

		cycles++
		if cycles%2000 == 0 {
			log.Println(num, "2k cycle. cnt", cnt, "r", ready_cnt, "b", bad_cnt, "c", cycles, time.Since(start))
			s.totalData += int64(last_cnt)
			last_cnt = 0
		}

		if send_close {
			break
		}
	}
	conn.Close()
	s.totalData += int64(last_cnt)
}

func (s *Server) serve() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Println("started")
	num := 0

	s.checkStoppingFlag()

	for {
		conn, err := ln.Accept()
		s.tcpAccepts++
		//		log.Println("accepted", accepts, err)

		if err != nil {
			log.Println(err.Error())
			return
		}
		go s.handleConnection(conn, num)
		num++
	}
}

var prev_total int64 = 0

func dumpStat(s *Server) {
	util.DumpMemStat()
	log.Printf("%s a: %d, c: %d, gon: %d, total: %d (%d)\n",
		time.Now(),
		s.connAllocs, s.connCached, runtime.NumGoroutine(), s.totalData, (s.totalData-prev_total)/5)

	log.Println("accepted:", s.tcpAccepts, "released", s.connReleases, "stop", s.stopping, "ready", s.totalReady)
	prev_total = s.totalData
	s.checkStoppingFlag()
}

func dumpStatProc(s *Server) {
	ticker := time.NewTicker(time.Second * 5)
	for _ = range ticker.C {
		dumpStat(s)
	}
}

func main() {
	util.InitLog()
	runtime.GOMAXPROCS(runtime.NumCPU())

	s := newServer(81)
	go dumpStatProc(s)
	s.serve()
}
