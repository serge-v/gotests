// test console http server

package main

import (
    "fmt"
    "bufio"
    "os"
    "io"
    "net"
    "net/http"
    "log"
    "time"
    "reflect"
)

type RootHandler struct{
	count int
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.count++
//	log.Println(r.Method, count2)
	if (h.count % 1000 == 0) {
		log.Println("requests: ", h.count)
	}
//	dump_httpreq(r)
	fmt.Fprintln(w, "root handler:", h.count)
}

type HelpHandler struct{
	count int
}

func (h *HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	fmt.Fprintln(w, "help handler")
}

func dump_httpreq(req *http.Request) {
	st := reflect.TypeOf(req)
	fmt.Println(st)

	val := reflect.ValueOf(req).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		fmt.Printf("    %s: %v\n", typeField.Name, valueField.Interface())
	}
	
	fmt.Println("    === headers ===")
	for k, v := range req.Header {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func server2(conf Config) {

	var h1 RootHandler
	var h2 HelpHandler
	http.Handle("/", &h1)
	http.Handle("/help/", &h2)

	log.Println("starting on :", conf.Port)
	
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

type Config struct {
	Port int
}

func main() {

	conf := Config{
		Port: 4000,
		}
		
	w := io.MultiWriter(os.Stderr)

	log.SetOutput(w)
	
	server2(conf)
}
