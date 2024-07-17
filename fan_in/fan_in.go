package main

import (
	"fmt"
	"sync"
)


func main() {
	numProducers := 10

	inChans := make([]chan int, numProducers, numProducers)
	for i := range inChans {
		inChans[i] = make(chan int)
		go producer(inChans[i])
	}

	var count int
	var partialSum int
	batchSize := 20
	for val := range multiWorkerMerge(inChans) {
		if count == batchSize - 1 {
			fmt.Println(partialSum + val)
			count = 0
			partialSum = 0
		} else {
			count++
			partialSum += val
		}
	}

	fmt.Println(partialSum)
}

func producer(inChan chan<- int) {
	defer close(inChan)
	
	for i := 0; i < 10; i++ {
		inChan <- i
	}
}

func multiWorkerMerge(cs []chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	mergeWorker := func(c <-chan int) {
		defer wg.Done()
		for num := range c {
			out <- num
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go mergeWorker(c)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
