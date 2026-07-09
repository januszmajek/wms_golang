package main

import (
	"fmt"
	"go-tour-exercises/fibonacci"
)

func main() {
	f := fibonacci.Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

// func main() {
// 	fmt.Println(Sqrt(2))
// }
