package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Sum struct {
	mu  sync.Mutex
	sum int
}

func main() {
	var arg_func string
	if len(os.Args) > 1 {
		arg_func = os.Args[1]
		fmt.Printf("arg_func: %s\n", arg_func)
	}
	//fmt.Println("Hello, World!")
	if arg_func == "concurrency" {
		wg := sync.WaitGroup{}
		done_chan := make(chan bool)
		wait_chan := make(chan int)
		for i := 0; i < 10; i++ {
			fmt.Println("Starting")
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				//fmt.Printf("Starting %d\n", i)
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				wait_chan <- i
			}(i)
		}
		go func() {
			wg.Wait()
			close(done_chan)
			close(wait_chan)
		}()
		for {
			select {
			case i, ok := <-wait_chan:
				if !ok {
					fmt.Println("Exiting")
					return
				}
				fmt.Printf("Finished %d\n", i)
			case <-done_chan:
				fmt.Println("All Done :)")
				os.Exit(0)
			}
		}
	} else if arg_func == "mutex" {
		fmt.Println("Mutex")
		wg := sync.WaitGroup{}
		sum := Sum{}
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				sum.mu.Lock()
				fmt.Printf("Adding %d to sum of %d\n", i, sum.sum)
				sum.sum += i
				sum.mu.Unlock()
			}(i)
		}
		wg.Wait()
	} else {
		x := 10000000
		y := &x
		fmt.Printf("x memory address: %p\n", y)
	}
}
