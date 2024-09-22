package main

import (
	"fmt"
	"sync"
)

func genTasks(limit int, taskWg *sync.WaitGroup) []chan int {
	perChannelLimit := limit / 5
	var taskChannels []chan int
	for t := 1; t < 6; t++ {
		taskWg.Add(1)
		taskChannel := make(chan int)
		go func(multiplier int) {
			defer taskWg.Done()
			for i := (multiplier - 1) * perChannelLimit; i < perChannelLimit*multiplier; i++ {
				taskChannel <- i
			}
			defer close(taskChannel)
		}(t)
		taskChannels = append(taskChannels, taskChannel)
	}
	return taskChannels
}

func main() {
	limit := 20000
	taskWg := &sync.WaitGroup{}
	taskChannels := genTasks(limit, taskWg)
	doneChan := make(chan bool)
	go func() {
		taskWg.Wait()
		close(doneChan)
	}()
	workerWg := sync.WaitGroup{}
	for itr := 1; itr < 6; itr++ {
		workerWg.Add(1)
		go func(i int) {
			defer workerWg.Done()
			for {
				remainingChannels := 0
				for _, taskChannel := range taskChannels {
					chi, ok := <-taskChannel
					if !ok {
						continue
					}
					remainingChannels++
					processedResult := chi * chi
					fmt.Printf("Worker %d processed task #%d with result %d\n", i, chi, processedResult)
				}
				select {
				case <-doneChan:
					return
				default:
				}
			}
		}(itr)
	}
	workerWg.Wait()
}
