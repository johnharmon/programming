package main

import "io"

func ReadInput(input io.Reader, out chan byte) {
	b := make([]byte, 1)
	for {
		input.Read(b)
		out <- b[0]
	}
}
