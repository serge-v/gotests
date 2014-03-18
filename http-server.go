// test console http server

package main

import (
    "fmt"
    "bufio"
    "os"
    "net"
    "net/http"
    "log"
    "time"
)

type RootHandler struct{}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "root handler")
}

type HelpHandler struct{}

func (h HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "help handler")
}

func server1() {
	var h RootHandler
	http.ListenAndServe("localhost:4000", h)
}

func server2() {
	var h1 RootHandler
	var h2 HelpHandler

	http.HandleFunc("/", h1.ServeHTTP)
	http.HandleFunc("/help/", h2.ServeHTTP)

	s := &http.Server{
		Addr:           ":4001",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		}

	l, e := net.Listen("tcp", ":4001")
	if e != nil {
		log.Panicf(e.Error())
	}
	go s.Serve(l)
	
	fmt.Println("For test run: curl -v http://localhost:4001/help/")
	fmt.Println("Press ENTER to stop");
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func main() {
	server2()
}
