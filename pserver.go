package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

type PairedServerHandler struct {
}

type PairedServer struct {
	listener net.Listener
	server   *http.Server
	number   int
	endpoint string
	active   bool
}

func newServer(port int) *PairedServer {
	var handler PairedServerHandler

	s := new(PairedServer)
	s.endpoint = fmt.Sprintf(":%d", port)
	s.server = &http.Server{
		Addr:           s.endpoint,
		Handler:        &handler,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1024,
		StopChan:       make(chan bool, 1),
		StoppedChan:    make(chan bool),
	}
	return s
}

func (srv *PairedServer) bindAndGoServe() {
	retries := 0
	var err error
	for {
		srv.listener, err = net.Listen("tcp", srv.endpoint)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 1)
		retries++
		if retries%2000 == 0 {
			log.Println("binding", err.Error())
		}
		if retries > 10000 {
			log.Panicln("cannot bind")
		}
	}

	go srv.server.Serve(srv.listener)
	srv.active = true
	log.Println("bound. retries: ", retries)
}

func (srv *PairedServer) activateSibling() {
	log.Println("activating sibling")
	url := fmt.Sprintf("http://127.0.0.1:%d/c/e", siblingPort)
	_, err := http.Get(url)
	if err != nil {
		log.Println("command error. continue active mode.", err)
		return
	}
	log.Println("activate command sent")
	//	time.Sleep(time.Millisecond*5)
	commandChan <- Cstop
}

func (srv *PairedServer) stop() {
	srv.server.StopChan <- true
	srv.listener.Close()
	<-srv.server.StoppedChan
	srv.active = false
	//	time.Sleep(time.Second * 1) // wait to complete go handlers
	//	log.Println("server stopped")
}

func dumpStat(instant uint64, active bool) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("mem: %dMb, a: %t, instant: %d\n", m.Sys/1024/1024, active, instant)
}

var stats = debug.GCStats{
	PauseQuantiles: make([]time.Duration, 5, 5),
	Pause:          make([]time.Duration, 1, 1),
}
var prevms runtime.MemStats
var ms runtime.MemStats

func collectGarbage() {
	time.Sleep(time.Second * 5)
	runtime.ReadMemStats(&ms)
	debug.SetGCPercent(100)
	runtime.GC()
	debug.FreeOSMemory()
	debug.ReadGCStats(&stats)

	fmt.Printf("GC pause: %dms, mallocs: %d, frees: %d\n",
		stats.Pause[0]/1000000,
		ms.Mallocs-prevms.Mallocs,
		ms.Frees-prevms.Frees)

	debug.SetGCPercent(-1)
	runtime.ReadMemStats(&prevms)
}

// === control server ===

type Command int

const (
	Cstart Command = iota
	CstartSibling
	Cstop
)

var commandChan = make(chan Command, 10)

type ControlHandler struct {
}

func controlServer(endpoint string) error {
	var h ControlHandler
	s := &http.Server{
		Addr:           endpoint,
		Handler:        &h,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}

	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		return err
	}

	go s.Serve(l)
	log.Println("control server started on", endpoint)
	return nil
}

func (s *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.URL.Path == "/c/e" {
		commandChan <- Cstart
	} else if r.URL.Path == "/c/d" {
		commandChan <- CstartSibling
	}
	fmt.Fprintln(w, r.URL.Path)
}
