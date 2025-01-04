package main

import (
	"fmt"
	"sync"
)

func genTasks(limit int) (taskChannel chan int) {
	taskChannel = make(chan int)
	go func() {
		for i := 0; i < limit; i++ {
			taskChannel <- i
		}
		defer close(taskChannel)
	}()
	return taskChannel
}
func main() {
	limit := 20000
	taskChannel := genTasks(limit)

	workerWg := sync.WaitGroup{}
	for itr := 1; itr < 6; itr++ {
		workerWg.Add(1)
		go func(i int) {
			defer workerWg.Done()
			//workerResults := []int{}
			for chi := range taskChannel {
				processedResult := chi * chi
				//workerResults = append(workerResults, processedResult)
				fmt.Printf("Worker %d processed task #%d with result %d\n", i, chi, processedResult)
			}
		}(itr)
	}
	workerWg.Wait()
}
