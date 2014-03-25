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
	"sync"
	"./version"
	"bufio"
)

const N = 10000
const SENDERS = 10

type RootHandler struct {
}

var userrec = make(chan string)
var done = make(chan int)

var appLocker sync.RWMutex

type App struct {
	Codes map[string]string
}

var app = &App{
	Codes: make(map [string]string, 1000),
}

func getApp() *App {
	appLocker.RLock()
	defer appLocker.RUnlock()
	return app
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	log.Println("req")
/*	if h.count%1000 == 0 {
		log.Println("requests: ", h.count)
	}*/
	
	a := getApp()
	userrec <- fmt.Sprintf("%s, %s", a.Codes["test"], r.URL)
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
	
	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://localhost:4001/sadasdasd/dfgdfgdfgdfsg/dfgdfgdsfg/dfgsdfgdfsg/dfgdfsgdsfgdsfg/dfgdfgdfg", nil)
	
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
	
	for i := 0; i < SENDERS; i++ {
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
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l, e := net.Listen("tcp", ":8080")
	if e != nil {
		log.Panicf(e.Error())
	}

	go s.Serve(l)

	start := time.Now()

	file, err := os.OpenFile("2.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}

	for {
		select {
			case msg := <-userrec:
				fmt.Fprintln(file, msg)
				
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("sender: elapsed: %v, speed: %.1f kps\n", elapsed, N*SENDERS/elapsed.Seconds()/1000)
	fmt.Printf("app: %+v\n", app)
	
	fmt.Println("Press ENTER to stop")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}

type Config struct {
	Port int
}

func updateCache() {
	ticker := time.NewTicker(time.Second)

	for t := range ticker.C {
		m1 := make(map[string]string)
		m1["test"] = fmt.Sprintf("test %v", time.Now())
		if false {
			fmt.Println(t)
		}

		a := App{
			Codes: m1,
		}
		
	
		appLocker.Lock()
		app = &a
		appLocker.Unlock()
	}

}

func main() {
	fmt.Println("version:", version.HEAD)

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

	go updateCache()
	server2(conf)
}
