package main

import (
	"fmt"
	"sync"
)

func genTasks(limit int) ([]chan int, *sync.WaitGroup) {
	perChannelLimit := limit / 5
	taskWg := &sync.WaitGroup{}
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
	return taskChannels, taskWg
}
func main() {
	limit := 20000
	taskChannels, taskWg := genTasks(limit)
	doneChan := make(chan bool)
	go func() {
		taskWg.Wait()
		fmt.Println("all task generators done")
		close(doneChan)
	}()
	workerWg := sync.WaitGroup{}
	for itr := 1; itr < 6; itr++ {
		workerWg.Add(1)
		go func(i int) {
			defer workerWg.Done()
			//workerResults := []int{}
			for {
				//allClosed := true
				remainingChannels := 0
				for _, taskChannel := range taskChannels {
					chi, ok := <-taskChannel
					if !ok {
						continue
					}
					//allClosed = false
					remainingChannels++
					processedResult := chi * chi
					fmt.Printf("Worker %d processed task #%d with result %d\n", i, chi, processedResult)
				}
				//				if allClosed {
				//					break
				//				}
				//				if allClosed {
				//					return
				//				}
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
