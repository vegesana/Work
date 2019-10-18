package main

import "fmt"

func main() {
	fn := Hello
	fn(10)
}

func Hello(val int) int {
	fmt.Println("value input is ", val)
	return val
}
