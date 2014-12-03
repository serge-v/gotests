package main

import "fmt"

type node struct {
	next  *node
	value string
}

type linkedList struct {
	first *node
	last *node
}

func (l *linkedList) add(v string) {
	n := new(node)
	n.value = v

	if l.first == nil {
		l.first = n
	} else {
		l.last.next = n
	}

	l.last = n
}

func (l *linkedList) insert(v string) {
	n := new(node)
	n.value = v

	if l.first == nil {
		l.last = n
	}

	n.next = l.first
	l.first = n
}

func (l *linkedList) reversed() *linkedList {
	nl := linkedList{}
	for n := l.first; n != nil; n = n.next {
		nl.insert(n.value)
	}
	return &nl
}

func (l *linkedList) print() {
	for n := l.first; n != nil; n = n.next {
		fmt.Print(n.value, " ")
	}
	fmt.Println()
}

func main() {

	list := new(linkedList)

	list.add("1")
	list.add("2")
	list.add("3")
	list.add("4")
	list.add("5")

	list.print()
	list.reversed().print()
}
