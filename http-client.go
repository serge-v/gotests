// test http client

package main

import (
	"fmt"
	"net/http"
//	"reflect"
	"time"
	"runtime"
)

const N = 10000
const SENDERS = 10

func dump_httpresp(resp *http.Response) {
	fmt.Println("status: ", resp.Status)
/*
	st := reflect.TypeOf(resp)
	fmt.Println(st)

	val := reflect.ValueOf(resp).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		typeName := valueField.Type().Name()

		fmt.Printf("%s(%s),\t\t\t Field Value: %v\n", typeField.Name, typeName, valueField.Interface())
	}

	fmt.Println("=== headers ===")
	for k, v := range resp.Header {
		fmt.Printf("%20s = %20s\n", k, v)
	}*/
}

func sender(num int, done chan int) {

	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	
	client := &http.Client{Transport: tr}
	
	url := "http://10.68.20.200:8080/data_provider/appnexus?uid=12000000000&aid=11000000000&country=US&seller=15000&url=http%3A%2F%2Fwww.test.com%2F"
	req, _ := http.NewRequest("GET", url, nil)
	
	for i := 0; i < N; i++ {
		
		resp, err := client.Do(req)
		
		if err != nil {
			fmt.Println(err)
			break
		}

		defer resp.Body.Close()
	}
	done <- num
}

func main() {
	runtime.GOMAXPROCS(4)

	start := time.Now()

	var done = make(chan int)
	
	for i := 0; i < SENDERS; i++ {
		go sender(i+1, done)
	}
	
	cnt := 0

	for num := range (done) {
		cnt++
		fmt.Println(num, "done")
		if cnt == SENDERS {
			break
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("sender: elapsed: %v, speed: %.1f kps\n", elapsed, N*SENDERS/elapsed.Seconds()/1000)
}
