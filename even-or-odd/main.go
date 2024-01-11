package main

import "fmt"

func main() {
	numbers := []int{5, 6, 7, 8, 1, 2, 3, 4}

	for n := range numbers {
		if n%2 == 0 {
			fmt.Printf("Number %d is even\n", n)
		} else {
			fmt.Printf("Number %d is odd\n", n)
		}
	}
}
