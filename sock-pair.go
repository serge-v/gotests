package main

import (
	"fmt"
	"log"
//	"net"
	"net/http"
	"os"
	"time"
	"bufio"
//	"runtime"
//	"bufio"
	"runtime/debug"
//	"runtime/pprof"
//	"strconv"
	"strings"
)

var siblingPort int

func (s* PairedServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintln(w, "resp")
}

var reader = bufio.NewReader(os.Stdin)

func readCommands() {
	for {
		s, _ := reader.ReadString('\n')
		commandChan <- strings.Trim(s, "\n")
	}
}

func loop() {

	go readCommands()

	var s *PairedServer
//	var err error

	ticker := time.NewTicker(time.Second*2)
	enabled_ticks := 0

	for {
		select {
		case <-ticker.C:
			dumpStat(0, s != nil && s.active)
			if s != nil && s.active {
				enabled_ticks++
				if enabled_ticks > 10 {
					commandChan <- "disable"
				}
			}

		case cmd := <-commandChan:
			log.Println("cmd:", cmd)
			if cmd == "q" {
				return
			} else if cmd  == "disable" {
				if s == nil {
					log.Println("server already stopped. start first")
					continue
				}

				log.Println("enable sibling")
				if !s.activateSibling() {
					continue
				}
				log.Println("enable command send")
				time.Sleep(time.Millisecond*5)
				s.stop()
				s = nil
				enabled_ticks = 0
				time.Sleep(time.Second * 5)
				log.Println("server stopped")
				collectGarbage()
			} else if cmd == "enable" {
				if s != nil {
					log.Println("server already started. stop first")
					continue
				}
				s = newServer(8081)
				s.bindAndGoServe()
			}
		}
	}
}

func main() {
	var err error

	debug.SetGCPercent(-1)

	err = controlServer(":8082")
	siblingPort = 8083

	if err != nil {
		err = controlServer(":8083")
		siblingPort = 8082
	} else {
		commandChan <- "enable"
	}

	loop()
}
