// test console http server

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"time"
	"runtime"
)

type RootHandler struct {
}

var userrec = make(chan string)
var done = make(chan int)


func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	log.Println("req")
	userrec <- "some useful info"
/*	if h.count%1000 == 0 {
		log.Println("requests: ", h.count)
	}*/
}

type HelpHandler struct {
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

func singleSender(num int) {
	const N = 10000
	
	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://localhost:4001/sadasdasd/", nil)
	
	for i := 0; i < N; i++ {
		
		resp, err := client.Do(req)
		
		if err != nil {
			fmt.Println(err)
			break
		}
		
/*		if i == N-1 {
			dump_httpresp(resp)
		}
*/		defer resp.Body.Close()
	}
//	fmt.Println("sender: sent: ", num)
	done <- num
}

func sender() {
	
	for i := 0; i < 10; i++ {
		go singleSender(i+1)
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
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l, e := net.Listen("tcp", ":4001")
	if e != nil {
		log.Panicf(e.Error())
	}

	go s.Serve(l)
	go sender()

	start := time.Now()

	file, err := os.OpenFile("2.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}

	donecount := 0

	for {
		select {
			case msg := <-userrec:
				fmt.Fprintln(file, msg)
			case num := <-done:
				donecount++
				fmt.Println("done:", num)
				
		}
		
		if donecount == 10 {
			break
		}
	}

	elapsed := time.Since(start)
	fmt.Println("sender: elapsed:", elapsed)
/*	
	fmt.Println("Press ENTER to stop")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
*/
}

type Config struct {
	Port int
}

func main() {
	file, err := os.OpenFile("1.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}
	log.SetOutput(file)
	log.Println("started")

	runtime.GOMAXPROCS(4)

	conf := Config{
		Port: 4000,
	}

	server2(conf)
}
