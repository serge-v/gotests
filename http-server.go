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
	"strconv"
	"strings"
)

type RootHandler struct {
}

const defInterval = 10000
var gcInterval = time.Duration(time.Millisecond * defInterval)
var ticker = time.NewTicker(time.Millisecond * defInterval)
var reloadTicker = make(chan bool)


func (s* RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	err = r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "parse error\n")
		return
	}

	if r.URL.Path == "/set/gc" {
		if len(r.Form["interval"]) > 0 {
			val, err := strconv.Atoi(r.Form["interval"][0])
			if err != nil {
				fmt.Fprintf(w, "parameter error\n")
				return
			}
			gcInterval = time.Millisecond * time.Duration(val)
			fmt.Fprintf(w, "new gc.inverval is %s\n", gcInterval)
			reloadTicker <- true
		} else {
			fmt.Fprintf(w, "no parameters specified\n")
		}
	} else if strings.HasPrefix(r.URL.Path, "/set/") {
		fmt.Fprintf(w, "%+v\n", r)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Keep-Alive", "timeout=120, max=50000")
		fmt.Fprintf(w, "resp\n")
	}
}

var stats = debug.GCStats{
	PauseQuantiles: make([]time.Duration, 5, 5),
	Pause:          make([]time.Duration, 1, 1),
}

var prevms runtime.MemStats
var ms runtime.MemStats

func collectGarbage() {

	runtime.ReadMemStats(&ms)

	debug.SetGCPercent(100)
	runtime.GC()
	debug.ReadGCStats(&stats)

	fmt.Printf("pause: %dms, mallocs: %d, frees: %d\n", stats.Pause[0] / 1000000, ms.Mallocs - prevms.Mallocs, ms.Frees - prevms.Frees)

	runtime.ReadMemStats(&prevms)
	
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
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
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
	
	go func() {
		reader.ReadString('\n')
		quit <- true
	}()

	freeMemory := false

	loop:
	for {
		select {
		case <- ticker.C:
			dumpHeap()
			collectGarbage()
			if freeMemory {
				debug.FreeOSMemory()
				log.Println("FreeOSMemory called")
				freeMemory = false
			}
		case <- quit:
			break loop
		case <- reloadTicker:
			ticker.Stop()
			ticker = time.NewTicker(gcInterval)
			log.Printf("new gc.inverval is %s\n", gcInterval)
			freeMemory = true
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
