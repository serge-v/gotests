package main

import (
	"io"
	"log"
	"os"
)

type coloredWriter struct {
	w     io.Writer
	debug bool
}

func (cw *coloredWriter) Write(p []byte) (n int, err error) {
	if len(p) < 21 {
		return cw.w.Write(p)
	}

	if p[20] == 'E' {
		cw.w.Write([]byte("\x1b[31m"))
		n, err = cw.w.Write(p)
		cw.w.Write([]byte("\x1b[0m"))
	} else if p[20] == 'W' {
		cw.w.Write([]byte("\x1b[35m"))
		n, err = cw.w.Write(p)
		cw.w.Write([]byte("\x1b[0m"))
	} else if p[20] == 'I' {
		cw.w.Write([]byte("\x1b[35m"))
		n, err = cw.w.Write(p)
		cw.w.Write([]byte("\x1b[0m"))
	} else if p[20] == 'D' {
		cw.w.Write([]byte("\x1b[36m"))
		n, err = cw.w.Write(p)
		cw.w.Write([]byte("\x1b[0m"))
	} else {
		n, err = cw.w.Write(p)
	}
	return
}

func setLog() {
	cw := &coloredWriter{w: os.Stderr, debug: true}
	log.SetOutput(cw)
}

func main() {
	setLog()
	log.Println("ERROR", "this is error")
	log.Println("WARN", "this is warning")
	log.Println("ERROR", "this is error")
	log.Println("WARN", "this is warning")
	log.Println("INFO", "this is warning")
	log.Println("DEBUG", "this is debug")
	log.Println("INFO", "debug was filtered")
}
