package main

import (
	"os"
	"fmt"
	"net/url"
	"runtime"
	"runtime/pprof"
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

func main() {

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
