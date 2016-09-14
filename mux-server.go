package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"io"
	"fmt"
)

func dumpReq(req *http.Request, w io.Writer) {
	fmt.Fprintf(w, "%+v\n", req)

	buf, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(buf))
	println(string(buf));
}

type dumpHandler struct {
}

func (dumpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	dumpReq(r, rw)
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/dump", &dumpHandler{})
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
