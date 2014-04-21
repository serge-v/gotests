package main

import (
	"os"
	"fmt"
	"net/url"
	"runtime"
	"runtime/pprof"
	"bytes"
)

func dumpHeapProfile() {
	p := pprof.Lookup("heap")
	if err := p.WriteTo(os.Stdout, 1); err != nil {
		fmt.Println("heap:", err.Error())
	}
}

var m1 runtime.MemStats
var m2 runtime.MemStats

func dumpMemStat() {
	runtime.ReadMemStats(&m2)
	fmt.Printf("M %d %d (%d), F %d %d (%d), gr: %d\n",
		m1.Mallocs, m2.Mallocs, m2.Mallocs - m1.Mallocs,
		m1.Frees, m2.Frees, m2.Frees - m1.Frees,
		runtime.NumGoroutine())
	runtime.ReadMemStats(&m1)
}

func test3() {
	urls := [][]byte{
		[]byte("http://www.test.com:81?aid=5655001397149476941&ip=123.123.128.75&uid=223411434234234&country=US"),
		[]byte("http://www.test.com:81?aid=5655001397149476941&ip=134.123.128.75 "),
		[]byte("http://www.test.com:81?aid=5655001397149476941&ip=123.123.128.75"),
		[]byte("http://www.test.com:81?aid=5655001397149476941&ip="),
		[]byte("http://www.test.com:81?aid=5655001397149476941"),
	}

	for _, u := range(urls) {
		idx1 := bytes.Index(u, []byte("ip="))
		idx2 := bytes.IndexAny(u[idx1+3:], "& ")
		if idx1 == -1 {
			continue
		}
		if idx2 == -1 {
			idx2 = len(u[idx1+3:])
		}
		fmt.Println("idx2", idx2)
		fmt.Println("idx1", idx1)
		fmt.Printf("ip: '%s'\n", string(u[idx1+3:idx1+3+idx2]))
	}
}

func test2() {
	rawurl := "http://www.test.com:81?test=aaa&test1=bbb&test3=ccc"

	for i := 0; i < 10000; i++ {
		_, err := url.Parse(rawurl)
		if err != nil {
			fmt.Println(err)
			break
		}

//		fmt.Fprintln(os.Stdout, "test")
//		fmt.Fprintln(os.Stdout, "test2")
	}
	dumpHeapProfile()
	dumpMemStat()
}

func main() {
	test3()
}
