/*
package main

import (

	"fmt"
	"math"

)

	func main() {
		investment_amount := 1000 // initial investment amount
		returnrate := 0.107       // return rate
		years := 30
		result := float64(investment_amount) * math.Pow((1+returnrate), float64(years))
		fmt.Printf("The investment amount after %d years is: %f", years, result)
	}
*/
package main

import "fmt"

func main() {
	// Integer formatting
	fmt.Printf("%d\n", 123) // Decimal
	fmt.Printf("%b\n", 123) // Binary
	fmt.Printf("%o\n", 123) // Octal
	fmt.Printf("%x\n", 123) // Hexadecimal
	fmt.Printf("%X\n", 123) // Uppercase Hexadecimal

	// Floating-point formatting
	fmt.Printf("%f\n", 123.456) // Decimal
	fmt.Printf("%e\n", 123.456) // Scientific notation
	fmt.Printf("%g\n", 123.456) // Compact representation

	// String formatting
	fmt.Printf("%s\n", "hello") // String
	fmt.Printf("%q\n", "hello") // Double-quoted string
	fmt.Printf("%x\n", "hello") // Hexadecimal

	// Width and precision
	fmt.Printf("|%6d|\n", 123)       // Width 6
	fmt.Printf("|%-6d|\n", 123)      // Width 6, left-justified
	fmt.Printf("|%6.2f|\n", 123.456) // Width 6, precision 2
	fmt.Printf("|%.2f|\n", 123.456)  // Precision 2

	// Struct formatting
	type Person struct {
		Name string
		Age  int
	}
	p := Person{"Alice", 30}
	fmt.Printf("%v\n", p)  // Default format
	fmt.Printf("%+v\n", p) // Field names
	fmt.Printf("%#v\n", p) // Go-syntax representation
}
