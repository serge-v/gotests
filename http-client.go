// test http client

package main

import (
	"fmt"
	"net/http"
	"time"
	"runtime"
	"io"
	"io/ioutil"
	"os"
)

const N = 10000
const SENDERS = 2000

type result struct {
	num int
	max time.Duration
	m int
	m0 int
	m5 int
	m10 int
	done bool
}

func sender(num int, done chan result) {

	tr := &http.Transport{
		DisableKeepAlives: false,
	}
	
	client := &http.Client{Transport: tr}
	var res result
	urlfmt := "http://%s:8080/data_provider/appnexus?uid=%d&ip=1.2.3.4"
	for i := 0; i < N; i++ {
		
		url := fmt.Sprintf(urlfmt, conf.Host, 12000000000+i)
		
		start := time.Now()
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		elapsed := time.Since(start)
		if elapsed > res.max {
			res.max = elapsed
		}
		
		if elapsed > time.Millisecond*10 {
			res.m10++
		} else if elapsed > time.Millisecond*5 {
			res.m5++
		} else {
			res.m0++
		}
		res.m++

		if err != nil {
			fmt.Println(1, err)
			os.Exit(1)
		}

		io.Copy(ioutil.Discard, resp.Body); // need to read body completely otherwize keep-alive doesn't work
		resp.Body.Close()

		if i%10000 == 0 {
			res.num = num
			done <- res
			res = result{}
		}
	}
	res.num = num
	res.done = true
	done <- res
}

func main() {

	if !parseConf() {
		return
	}

	runtime.GOMAXPROCS(4)

	start := time.Now()

	var done = make(chan result)
	
	for i := 0; i < SENDERS; i++ {
		go sender(i+1, done)
	}
	
	cnt := 0
	done_cnt := 0
	
	var total result

	for res := range (done) {
		if res.max > total.max {
			total.max = res.max
		}
		fmt.Printf("%d: max: %s, m: %d, m0: %d, m5: %d, m10: %d\n",
			res.num, res.max, res.m, res.m0, res.m5, res.m10)

		total.m10 += res.m10
		total.m5 += res.m5
		total.m0 += res.m0
		total.m += res.m
		cnt++
		if cnt%SENDERS == 0 {
			elapsed := time.Since(start)
			fmt.Printf("max: %s, m: %d, m0: %d, m5: %d, m10: %d, elapsed: %v, speed: %.1f kps\n",
				total.max, total.m, total.m0, total.m5, total.m10,
				elapsed, float64(total.m)/elapsed.Seconds()/1000)
			start = time.Now()
			total = result{}
		}
		if res.done {
			done_cnt++
			if done_cnt == SENDERS {
				break
			}
		}
	}
}
