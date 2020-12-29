package main

var test = []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}

func fib(n int) int {
	if n <= 0 {
		return 0
	}

	pprev := 0
	prev := 1
	println(prev)
	for i := 1; i < n; i++ {
		next := prev + pprev
		println(next)
		pprev = prev
		prev = next
	}
	return prev
}

func fibr(n int) int {
	if n == 1 || n == 2 {
		return 1
	}
	r := fibr(n-1) + fibr(n-2)
	return r
}

func main() {
	fib(50)
	println("---")
	for i := 2; i < 50; i++ {
		println(fibr(i))
	}
}
