package main

import (
	"runtime"
	"runtime/pprof"
	"os"
	"bufio"
	"fmt"
	"bytes"
	"log"
)

func dumpHeap() {
	p := pprof.Lookup("heap")
	if err := p.WriteTo(os.Stdout, 1); err != nil {
		fmt.Println("heap:", err.Error())
	}
}

var m1 runtime.MemStats
var m2 runtime.MemStats

var buff = make([]byte, 1000, 1000)


func main() {

	file, err := os.OpenFile("string-test~.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}
	
	fw := bufio.NewWriter(file)

	b := bytes.NewBuffer(buff)

//	dumpHeap()
	runtime.ReadMemStats(&m1)
//	s := ""

	for i := 0; i < 10000; i++ {
//		s = fmt.Sprintf("item number %d", i)
//		fmt.Println(s)
//		fmt.Fprintf(b, "item number %d", i)
//		fmt.Println(b)
		fmt.Fprintf(fw, "item number %d\n", i)
		b.Reset()
	}

	runtime.ReadMemStats(&m2)
//	dumpHeap()

	fmt.Printf("Mallocs %d %d (%d)\n", m1.Mallocs, m2.Mallocs, m2.Mallocs - m1.Mallocs)
	fmt.Printf("Frees   %d %d (%d)\n", m1.Frees, m2.Frees, m2.Frees - m1.Frees)
	fmt.Printf("b: %d\n", b.Len())
	
	fw.Flush()
}
