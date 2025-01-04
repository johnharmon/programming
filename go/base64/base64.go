package main

import (
	"fmt"
	"sync"
)

//func bitByteTest(byte) {
//	var x uint8 = 0
//	var n int = 2048
//
//	x << 8
//}

func goRoutine(n int, wg *sync.WaitGroup) (ch chan int) {
	c := make(chan int)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			c <- i
			defer wg.Done()
			fmt.Println("go routine: ", i)
		}(i)
	}
	return c
}

func main() {
	wg := sync.WaitGroup{}
	fmt.Println("created channel")
	ch := goRoutine(10, &wg)
	fmt.Println("go routine started")
	for range ch {
		chi := <-ch
		fmt.Printf("%d\n", chi)
	}
	fmt.Printf("Routines started, waiting on wg\n")
	wg.Wait()
	defer close(ch)
	//	args := os.Args[1:]
	//	s := strconv.Itoa(len(args))
	//	str_x, err := strconv.Atoi(s)
	//	if err != nil {
	//		fmt.Printf("Error: %s\n", err)
	//	}
	//
	// fmt.Printf("Converted: %d\n", str_x)
	// fmt.Println("placeholder")
}
