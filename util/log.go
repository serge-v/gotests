package util

import (
	"log"
	"fmt"
	"os"
	"io"
)

func InitLog() {

	logfile := "1.log"

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf(err.Error())
	}

	w := io.MultiWriter(os.Stderr, file)
	//	w := io.MultiWriter(file)
	log.SetOutput(w)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	fmt.Println("logging started to ", logfile)
}
