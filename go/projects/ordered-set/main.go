package main

import (
	"fmt"

	"harmonlab.io/set"
)

func main() {
	mySet := set.NewOrderedSet[int]()
	for i := 0; i < 5; i++ {
		res := mySet.Append(i)
		fmt.Printf("Change Applied: %t\n", res)
	}
	for i := 0; i < 5; i++ {
		res := mySet.Append(i)
		fmt.Printf("Change Applied: %t\n", res)
	}
}
