package main

import (
        "fmt"
        "net"
        "strings"
)

func getIp() string {
	var ip string
        ifs, _ := net.Interfaces()

        for _, iface := range(ifs) {
                addrs, _ := iface.Addrs()
                if len(addrs) > 0 {
			ip = addrs[0].String()
			ip = strings.Trim(ip, "[]")
			parts := strings.Split(ip, "/")
			ip = parts[0]
			if iface.Name != "lo" {
				return ip
			}
		}
        }
        
        return ip
}

func main() {
	fmt.Println(getIp())
}

