package main

import (
	"fmt"
	"sync"
)


func main() {
	var wg sync.WaitGroup

	numWorkers := 3
	inChan := make(chan int)
	outChan := make(chan int)

	go producer(inChan)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, inChan, outChan, &wg)
	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	for val := range outChan {
		fmt.Println(val)
	}
}

func producer(inChan chan<- int) {
	defer close(inChan)
	
	for i := 0; i < 10; i++ {
		inChan <- i
	}
}

func worker(i int, inChan <-chan int, outChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range inChan {
		outChan <- 2 * val + 1
	}
}
