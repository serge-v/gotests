package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"os"
	"io"
	"log"
)

var inputFile = flag.String("infile", "test.json", "Input file path")

type Message struct {
	Time int
	Userid int
	Node_ids []int
	Status int
	Key string
	Random string
	Enable_editing int
	Delay int
}

func main() {
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	dec := json.NewDecoder(file)
	
	var m Message

	if err := dec.Decode(&m); err == io.EOF {
		fmt.Println("EOF")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Time:           ", m.Time)
	fmt.Println("Userid:         ", m.Userid)
	fmt.Println("Node_ids:       ", m.Node_ids)
	fmt.Println("Status:         ", m.Status)
	fmt.Println("Key:            ", m.Key)
	fmt.Println("Random:         ", m.Random)
	fmt.Println("Enable_editing: ", m.Enable_editing)
	fmt.Println("Delay:          ", m.Delay)
}
