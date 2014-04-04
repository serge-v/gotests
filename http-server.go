// test console http server

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"runtime"
	"./version"
	"bufio"
	"runtime/debug"
	"runtime/pprof"
)

type RootHandler struct {
}

func (s* RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "resp\n")
}

func collectGarbage() {

	stats := debug.GCStats{
		PauseQuantiles: make([]time.Duration, 5, 5),
		Pause:          make([]time.Duration, 5, 5),
	}

	debug.SetGCPercent(100)
	runtime.GC()
	debug.ReadGCStats(&stats)
	log.Printf("stats after: %+v\n", stats)
	debug.SetGCPercent(-1)
}

func server(conf Config, daemon bool) {

	var h RootHandler
//	http.Handle("/", &h)

	endpoint := fmt.Sprintf(":%d", conf.Port)
	log.Println("starting on:", endpoint)

	s := &http.Server{
		Addr:           endpoint,
		Handler:        &h,
		ReadTimeout:    time.Second * 100,
		WriteTimeout:   time.Second * 100,
		MaxHeaderBytes: 1 << 20,
	}

	l, e := net.Listen("tcp", endpoint)
	if e != nil {
		log.Panicf(e.Error())
	}

	go s.Serve(l)

	quit := make(chan bool)

	fmt.Println("Press ENTER to stop")
	reader := bufio.NewReader(os.Stdin)
	
	t := time.NewTicker(time.Second * 10)

	go func() {
		reader.ReadString('\n')
		quit <- true
	}()

	loop:
	for {
		select {
		case <- t.C:
			dumpHeap()
			collectGarbage()
		case <- quit:
			break loop
				
		}
	}
}

var heapCount int

func dumpHeap() {
	f, _ := os.Create(fmt.Sprintf("heap~/heap%03d.prof", heapCount))
	heapCount++
	p := pprof.Lookup("heap")
	if err := p.WriteTo(f, 1); err != nil {
		log.Println("heap:", err.Error())
	}
	f.Close()
}

func main() {

	if !parseConf() {
		return
	}
	
	if conf.ShowVersion {
		fmt.Println("=== diff ===")
		fmt.Println(version.DIFF)
		fmt.Println("=== status ===")
		fmt.Println(version.STATUS)
		fmt.Println("=== head ===")
		fmt.Println(version.HEAD)
		return
	}

//	daemon(1, 1)

	fmt.Println("version:", version.HEAD)
/*
	file, err := os.OpenFile("1.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}

	log.SetOutput(file)
*/	log.Println("started")
	fmt.Println("log started")

	runtime.GOMAXPROCS(4)
	debug.SetGCPercent(-1)
	
	go func() {
 		log.Println(http.ListenAndServe("localhost:6060", nil))
 	}()
	
	server(conf, true)
}
