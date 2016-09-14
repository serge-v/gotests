package main

import (
	"fmt"
	"runtime/debug"
)

func div(divider int) (ret int) {

	ret = 1

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			b := debug.Stack()
			fmt.Println("=================\n", string(b))
			ret = 3
		}
	}()

	a := 1 / divider
	
	fmt.Println("a:", a)
	ret = 2
	return
}

func main() {
	fmt.Println("start")
	fmt.Println("div:", div(0))
	fmt.Println("stop")
}
