package main

import (
	"flag"
)

type Config struct {
	Host string
	Port int
	ShowVersion bool
}

var conf Config

func parseConf() bool {
	host := flag.String("host", "localhost", "Remote host")
	port := flag.Int("port", 8080, "Remote port")
	version := flag.Bool("version", false, "Show version, head hash and local source changes")
	flag.Parse()

	if !flag.Parsed() {
		return false
	}
	
	conf = Config{
		Host: *host,
		Port: *port,
		ShowVersion: *version,
	}

	return true
}
