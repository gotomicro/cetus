package main

import (
	"fmt"
)

func main() {
	a := 1 << 0
	b := 1 << 1
	c := 3

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	fmt.Println()
	fmt.Println(c ^ a)
	fmt.Println(c ^ b)

	fmt.Println()
	fmt.Println(c &^ a)
	fmt.Println(c &^ b)

	fmt.Println()
	fmt.Println(c | a)
	fmt.Println(c | b)
}
