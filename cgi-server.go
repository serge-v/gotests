package main

import (
	"net/http"
	"net/http/cgi"
	"fmt"
)

func main() {

	fmt.Println("starting server on http://localhost:9001")
	
	h := cgi.Handler{}
	h.Path = "http-cookie-server"
	
	panic(http.ListenAndServe(":9001", &h))
}
