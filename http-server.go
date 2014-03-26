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
)

const N = 10000
const SENDERS = 10

type RootHandler struct {
}

var userrec = make(chan string)
var done = make(chan int)

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("req:\n%v\n", r)))
	w.Write([]byte(fmt.Sprintf("url:\n%v\n", r.URL)))
	userrec <- fmt.Sprintf("url:\n%v\n", r.URL)
}

type HelpHandler struct {
	count int
}

func (h *HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	fmt.Fprintln(w, "help handler")
}

func server(conf Config) {

	var h1 RootHandler
	var h2 HelpHandler
	http.Handle("/", &h1)
	http.Handle("/help/", &h2)

	log.Println("starting on :", conf.Port)
	endpoint := fmt.Sprintf(":%d", conf.Port)

	s := &http.Server{
		Addr:           endpoint,
		Handler:        nil,
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
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

	loop:
	for {
		select {
			case msg := <-userrec:
				log.Println(msg)
			case <- quit:
				break loop
				
		}
	}
}

func main() {

	if !parseConf() {
		return
	}
	
	if *conf.ShowVersion {
		fmt.Println("=== diff ===")
		fmt.Println(version.DIFF)
		fmt.Println("=== status ===")
		fmt.Println(version.STATUS)
		fmt.Println("=== head ===")
		fmt.Println(version.HEAD)
		return
	}

	fmt.Println("version:", version.HEAD)

	file, err := os.OpenFile("1.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}
	log.SetOutput(file)
	log.Println("started")
	
	runtime.GOMAXPROCS(4)
	server(conf)
}
