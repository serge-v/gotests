package util

import (
	"log"
	"runtime"
	"runtime/pprof"
	"time"
	"os"
	"fmt"
)

var m1 runtime.MemStats
var m2 runtime.MemStats

func DumpMemStat() {
	runtime.ReadMemStats(&m2)
	log.Printf("%s mem: %d %d (%d), F %d %d (%d), gr: %d\n",
		time.Now(),
		m1.Mallocs, m2.Mallocs, m2.Mallocs-m1.Mallocs,
		m1.Frees, m2.Frees, m2.Frees-m1.Frees,
		runtime.NumGoroutine())
	runtime.ReadMemStats(&m1)
}

func DumpHeap() {
	p := pprof.Lookup("heap")
	if err := p.WriteTo(os.Stdout, 1); err != nil {
		fmt.Println("heap:", err.Error())
	}
}
