package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/cgi"
	"runtime"
)

var local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type handler struct {
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>It works!</p></body></html>")
}

func main() {
	flag.Parse()
	var err error
	
	h := &handler{}

	if *local != "" {
		err = http.ListenAndServe(*local, h)
	} else {
		err = cgi.Serve(h)
	}
	if err != nil {
		log.Fatal(err)
	}
}
