package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
)

func main() {

	c, err := ftp.Connect("localhost:21")
	if err != nil {
		fmt.Println(err)
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Println(err)
	}

	err = c.NoOp()
	if err != nil {
		fmt.Println(err)
	}
}
