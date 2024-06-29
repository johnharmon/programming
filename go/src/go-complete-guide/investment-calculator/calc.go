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
