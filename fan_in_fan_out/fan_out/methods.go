package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func (c *ChanFanOut) Batch() {
	wg := &sync.WaitGroup{}
	batchChan := make(chan []int64, c.workerCount)
	for i := 0; i < c.workerCount; i++ {
		batchChan <- c.nums[i * len(c.nums) / c.workerCount : (i + 1) * len(c.nums) / c.workerCount]
	}
	close(batchChan)


	for batch := range batchChan {
		wg.Add(1)
		go func(batch []int64) {
			defer wg.Done()
			var answer int64
			for _, num := range batch {
				answer += num
			}
			fmt.Println(c.message(answer))
		}(batch)
	}
	wg.Wait()
}

func (c *ChanFanOut) RoundRobin() {
	wg := &sync.WaitGroup{}
	workerChans := make([]chan int64, c.workerCount)
	for i := range workerChans {
		workerChans[i] = make(chan int64, len(c.nums) / c.workerCount + 1)
	}
	for i, num := range c.nums {
		workerChans[i % c.workerCount] <- num
	}
	for _, ch := range workerChans {
		close(ch)
		wg.Add(1)
		go func(workerChan chan int64) {
			defer wg.Done()
			var answer int64
			for num := range workerChan {
				answer += num
			}
			fmt.Println(c.message(answer))
		}(ch)
	}
	wg.Wait()
}

func (c *ChanFanOut) Random() {
	wg := &sync.WaitGroup{}
	workerChans := make([]chan int64, c.workerCount)
	reasonableBufferSize := 1000
	for i := range workerChans {
		workerChans[i] = make(chan int64, reasonableBufferSize)
	}
	for _, num := range c.nums {
		workerChans[rand.Int63n(int64(c.workerCount))] <- num
	}
	for _, ch := range workerChans {
		close(ch)
		wg.Add(1)
		go func(workerChan chan int64) {
			defer wg.Done()
			var answer int64
			for num := range workerChan {
				answer += num
			}
			fmt.Println(c.message(answer))
		}(ch)
	}
	wg.Wait()
}
