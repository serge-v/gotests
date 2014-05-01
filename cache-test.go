package main

import (
	"fmt"
	"sync"
	"time"
)

var m map[string]int
var mx sync.RWMutex

func getMap() map[string]int {
	//	mx.RLock()
	//	defer mx.RUnlock()
	return m
}

func reader(idx int) {
	var prev int
	for i := 0; i < 100000; i++ {
		m1 := getMap()
		curr := m1["version"]
		if curr != prev {
			fmt.Println(idx, ":", curr, ", i:", i)
			prev = curr
		}
		//		time.Sleep(time.Millisecond)
	}
	fmt.Println(idx, ": done")
}

func updater() {
}

func main() {
	m = make(map[string]int)
	m["version"] = 0

	ticker := time.NewTicker(time.Second)

	for i := 0; i < 10; i++ {
		go reader(i + 10)
	}

	ver := 1

	for t := range ticker.C {
		m1 := make(map[string]int)
		m1["version"] = ver
		if false {
			fmt.Println(t)
		}
		ver++

		mx.Lock()
		m = m1
		mx.Unlock()
	}
}
