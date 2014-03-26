package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Port int
	ShowUsage bool
	ShowVersion bool
}

var conf Config

func parseConf() bool {
	port := flag.Int("port", 8080, "Port to listen")
	usage := flag.Bool("help", false, "Show this help")
	version := flag.Bool("version", false, "Show version, head hash and local source changes")
	flag.Parse()

	if *usage || !flag.Parsed() {
		flag.VisitAll(func(f *flag.Flag){
			fmt.Println(f.Usage)
		})
		return false
	}
	
	conf = Config{
		Port: *port,
		ShowUsage: *usage,
		ShowVersion: *version,
	}

	return true
}
