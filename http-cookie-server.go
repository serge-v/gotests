package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"runtime"
)

var (
	local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
	date string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}


type handler struct {
}

func setCookies(w http.ResponseWriter) {
	c := &http.Cookie{}
	c.Name = "uid"
	c.Value = "100"
	http.SetCookie(w, c)
	c.Name = "uid2"
	c.Value = "200"
	http.SetCookie(w, c)
}

func getUid(cookies []*http.Cookie) (uid string, ok bool) {
	ok = false
	for _, c := range cookies {
		if c.Name == "uid" {
			uid = c.Value
			ok = true
			return 
		}
	}
	return
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	headers.Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "http-cookie-server: %s\n", date)
	uid, ok := getUid(r.Cookies())
	if ok {
		fmt.Fprintf(w, "uid: %s\n", uid)
	} else {
		setCookies(w)
	}
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
