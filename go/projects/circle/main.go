package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Coords struct {
	x int
	y int
}

func GetDistance(ax, ay, bx, by int) float64 {
	dx := ax - bx
	dy := ay - by
	return math.Sqrt(math.Pow(float64(dx), 2.0) + math.Pow(float64(dy), 2.0))
}

func LabelCells(radius int) ([][]bool, []Coords) {
	length, width := 2*radius+1, 2*radius+1
	square := make([][]bool, length)
	validCoords := make([]Coords, 0, length*width)
	validCoordsIndex := 0
	for row := range square {
		square[row] = make([]bool, width)
		for col := range square[row] {
			distToCenter := GetDistance(row, col, radius, radius)
			if int(math.Round(distToCenter)) < radius {
				square[row][col] = true
				validCoords = validCoords[0 : len(validCoords)+1]
				validCoords[validCoordsIndex] = Coords{x: col, y: row}
				validCoordsIndex++
			} else {
				square[row][col] = false
			}
		}
	}
	return square, validCoords //[0 : len(validCoords)-1]
}

func main() {
	// a := make([]byte, 2)
	square, coords := LabelCells(30)
	fmt.Print("\x1b[2J\x1b[0;0H")
	for row := range square {
		for _, b := range square[row] {
			if b {
				a := make([]byte, 2)
				a[0] = uint8(rand.Uint32()%0x5E) + 0x21
				a[1] = uint8(rand.Uint32()%0x5E) + 0x21
				fmt.Printf("%s", a)
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Print("\n")
	}
	for {
		coord := coords[rand.Uint32()%uint32(len(coords))]
		// if coord.x != 0 && coord.y != 0 {
		a := make([]byte, 2)
		a[0] = uint8(rand.Uint32()%0x5E) + 0x21
		a[1] = uint8(rand.Uint32()%0x5E) + 0x21
		fmt.Printf("\x1b[%d;%dH%s", coord.y+1, coord.x*2+1, a)
		time.Sleep(time.Millisecond * 2)
		//}

	}
}
