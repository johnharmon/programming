package main 

import (
	"fmt", 
	"os"
)

func factorial(n int) int {
	//target = os.Args[1] 
	if n == 0 {
		return 1
	} else {
		return n * factorial(n-1)
	}
}

func main() {
	target, err = os.Args[1]
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	result := factorial(target)
	fmt.Printf("Facttorial of %d is %d\n", target, result)
}
