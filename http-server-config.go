package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Port* int
	ShowUsage* bool
	ShowVersion* bool
}

var conf = Config{}

func parseConf() bool {
	conf.Port = flag.Int("port", 8080, "Port to listen")
	conf.ShowUsage = flag.Bool("help", false, "Show this help")
	conf.ShowVersion = flag.Bool("version", false, "Show version")
	flag.Parse()

	if !flag.Parsed() {
		flag.VisitAll(func(f *flag.Flag){
			fmt.Println(f.Usage)
		})
		return false
	}

	return true
}
