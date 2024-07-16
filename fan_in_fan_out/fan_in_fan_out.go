package main

import (
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)


type MySemaphore interface {
	Release()
	TryAcquire() bool
}

type UnweightedSemaphore struct {
	*semaphore.Weighted
	Capacity int64

	size int64
}

func (u *UnweightedSemaphore) Release() {
	u.Weighted.Release(1)
	u.size -= 1
}

func (u *UnweightedSemaphore) TryAcquire() bool {
	if u.size < u.Capacity && u.Weighted.TryAcquire(1) {
		u.size += 1
		return true
	}
	return false
}


func main() {
	var nums [1000]int64

	for i := range nums {
		nums[i] = 1
	}

	fmt.Println("Starting processing...")

	sem := UnweightedSemaphore{
		Weighted: semaphore.NewWeighted(7),
		Capacity: 8,
	}

	fmt.Println(useWorkerPool(&sem))
}

func useWorkerPool(sem MySemaphore) int64 {
	resultChan := make(chan int64)

	nums := make([]int64, 1000)

	i := 0

	var wg sync.WaitGroup
	for sem.TryAcquire() {
		fmt.Println(i)
		i += 1
		wg.Add(1)
		go func() {
			defer sem.Release()
			defer wg.Done()
			var result int64
			for _, num := range nums {
				result += num
			}
			resultChan <- result
		}()
	}

	wg.Wait()
	var answer int64
	for j := 0; j < i; j++ {
		result := <-resultChan
		answer += result
	}
	close(resultChan)
	return answer
}
