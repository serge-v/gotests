package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
//	"runtime"
//	"bufio"
//	"runtime/debug"
//	"runtime/pprof"
//	"strconv"
//	"strings"
)

type RootHandler struct {
}

func (s* RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func server() {

	var h RootHandler

	s := &http.Server{
		Addr:           ":8081",
		Handler:        &h,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}
	l, e := net.Listen("tcp", ":8081")
	if e != nil {
		log.Panicf(e.Error())
	}

	go s.Serve(l)

	fmt.Printf("Use Ctrl+C or 'kill -2 -p %d' to stop\n", os.Getpid())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	
	for {
		select {
		case <-signals:
			fmt.Println("stop")
			return
		}
	}
}

func main() {
	server()
}
