package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	cRed     = "\x1b[31m"
	cGreen   = "\x1b[32m"
	cBlue    = "\x1b[34m"
	cMagenta = "\x1b[35m"
	cNone    = "\x1b[0m"
)

func printMap(m interface{}, depth int) (wasObject bool) {
	indent := strings.Repeat("\t", depth)
	wasObject = false

	switch m.(type) {
	case map[string]interface{}:
		smap := m.(map[string]interface{})
		skeys := make([]string, 0, len(smap))
		for k := range smap {
			skeys = append(skeys, k)
		}
		sort.Strings(skeys)
		for _, k := range skeys {
			v := smap[k]
			if depth > 0 {
				fmt.Println()
			}
			fmt.Printf("%s%s:", indent, k)
			printMap(v, depth+1)
		}
		wasObject = true
	case []interface{}:
		fmt.Printf("[%d]", len(m.([]interface{})))
		for _, v := range m.([]interface{}) {
			if printMap(v, depth) {
				fmt.Printf("\n%s--", indent)
			}
		}
	case nil:
		fmt.Printf(" %snull%s", cMagenta, cNone)
	case bool:
		fmt.Printf(" %s%t%s", cMagenta, m, cNone)
	case string:
		fmt.Printf(" %s%s%s", cGreen, m, cNone)
	case json.Number:
		fmt.Printf(" %s%s%s", cMagenta, m, cNone)
	default:
		fmt.Printf("%s%s%s(%T)%s", indent, cRed, m, m, cNone)
	}

	if depth == 0 {
		fmt.Println()
	}

	return
}

func main() {
	f, _ := os.Open("2~.json")
	dec := json.NewDecoder(f)
	dec.UseNumber()
	var v map[string]interface{}
	err := dec.Decode(&v)
	if err != nil {
		panic(err)
	}

	printMap(v, 0)
}
