package main

import (
	"fmt"
)

type Engine struct {
	Volume int
}

type Wheels struct {
	Vendor string
}

type Car struct {
	*Wheels
	*Engine
	engVolume int
}

func NewEngine() *Engine {
	e := &Engine{
		5,
	}
	fmt.Println("NewEngine")
	return e
}

func NewWheels() *Wheels {
	w := &Wheels{
		"firestone",
	}
	fmt.Println("NewWheels")
	return w
}

func main() {

	car := Car{
		NewWheels(),
		NewEngine(),
		1,
	}

	fmt.Println(car)
}
