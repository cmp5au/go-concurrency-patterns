package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting...")
	squareChan := square()
	done := make(chan struct{}, 10)

	for i := 0; i < 10; i++ {
		go func(i int){
			squareChan <- i
			fmt.Println(<-squareChan)
			done <- struct{}{}
		}(i)
	}

	var count int
	for _ = range done {
		count++
		if count == 10 {
			break
		}
	}
	
	pingPong()
	fmt.Println("Done!")

	
}

// example of single-producer single-consumer: the function starts a server goroutine that writes to the same channel it reads from
// single-consumer because consumer B could be reading the request of consumer A off of the channel instead of the product of the producer
// single-producer because there's no point to having more if you can only have a single consumer
func square() chan int {
	ch := make(chan int)
	go func() {
		for x := range ch {
			ch <- x*x
		}
	}()
	return ch
}

// two goroutines sharing a channel for request/response-like comms 
func pingPong() {
	ch := make(chan int)

	player := func(ch chan int, msg string) {
		for {
			ct := <-ch
			fmt.Println(ct, msg)
			ch <- ct + 1
			if ct >= 9 {
				fmt.Println("Returning", msg)
				return
			}
		}
	}

	go player(ch, "ping")
	go player(ch, "pong")

	ch <- 0
	for {
		if i, ok := <-ch; ok {
			if i >= 11 {
				close(ch)
				break
			}
			ch <- i
		} else {
			close(ch)
			break
		}
	}
}

type Request struct {
	id string

}

type Response struct {

}
func singleProducerSingleConsumer() {
	// ch := make(chan Response)
}
