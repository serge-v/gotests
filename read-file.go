package main

import (
	"os"
	"io"
	"bufio"
	"fmt"
)

func readFile(fname string) (c chan string) {
	
	f, err := os.Open(fname)
	if err != nil {
		panic(err.Error())
	}

	c = make(chan string)

	r := bufio.NewReader(f)
	
	go func() {
		for {
			s, err := r.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("ERROR: ", err.Error())
				}
				break
			}
			c <- s
		}
		close(c)
		f.Close()
	}()
	return
}

func main() {
	fname := "IntegerArray.txt"
	
	count := 0
	tlen := 0
	
	for s := range readFile(fname) {
		count++
		tlen += len(s)
	}
	
	fmt.Println("lines:", count)
	fmt.Println("tlen: ", tlen)
	
}
