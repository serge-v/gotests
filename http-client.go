// test http client

package main

import (
	"fmt"
	"net/http"
//	"reflect"
	"time"
	"runtime"
)

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

func send(num int) {
	const N = 10000
	
	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", "http://localhost:4001/sadasdasd/", nil)
	fmt.Println("newreq", err)
	
	for i := 0; i < N; i++ {
		
		resp, err := client.Do(req)
		
		if err != nil {
			fmt.Println(err)
			break
		}

		if i%1000 == 0 {
			fmt.Println(i)
		}
		
		if i == N-1 {
			dump_httpresp(resp)
		}
		defer resp.Body.Close()
	}
	fmt.Println("sent", num)
}

func main() {
	runtime.GOMAXPROCS(4)

	start := time.Now()
	
	for i := 0; i < 20; i++ {
		go send(i+1)
	}
	
	send(0)


	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)
}
