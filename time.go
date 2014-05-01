package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()

	fmt.Println(now)

	fmt.Println(now.Add(time.Duration(1000 * time.Second)))
}
