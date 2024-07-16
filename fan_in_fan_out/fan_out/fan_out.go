package main

import (
	"fmt"
	"math/rand"
)

type FanOut interface {
	Batch()
	RoundRobin()
	Random()
}

type ChanFanOut struct {
	nums []int64
	workerCount int
	message func(int64) string
}


func main() {
	numSize := 14000
	cfo := ChanFanOut{
		nums: make([]int64, numSize),
		workerCount: 7,
		message: func(answer int64) string {
			return fmt.Sprintf("goroutine produced %v\n", answer)
		},
	}

	for i := range cfo.nums {
		cfo.nums[i] = rand.Int63n(11)
	}

	fmt.Println("Starting processing via fan out...")

	cfo.Batch()

	cfo.workerCount = 2

	cfo.RoundRobin()

	cfo.workerCount = 20

	cfo.Random()

	fmt.Println("Done!")
}
