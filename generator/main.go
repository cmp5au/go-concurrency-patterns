package main

import (
	"fmt"
	"math/rand"
	"time"
)

// The main idea behind a generator is to spawn some workers and return the channel that the workers use.

// This first generator is pretty much the same as Rob Pike's in the talk with a single worker goroutine
// that functions as a producer for the channel.
func boringGenerator(message string) <-chan string {
	boringChan := make(chan string)
	go func() {
		for i := 0; ; i++ {
			boringChan <- fmt.Sprintf("%v %s", i + 1, message)
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
		}
	}()
	return boringChan
}

// This generator is more of a request-response pattern where the channel is bidirectional. This type
// of generator only supports a single client and a single server because they must be synchronized.
func squareGenerator() chan int {
	squareChan := make(chan int)
	go func() {
		defer close(squareChan)
		for num := range squareChan {
			squareChan <- num * num
		}
	}()
	return squareChan
}

// This generator is to demonstrate the pattern of a quit channel being used for cleanup of workers that
// are no longer needed spawned by a longer-running main thread.
func quittableGenerator(quit chan struct{}) chan string {
	quittableChan := make(chan string)
	go func() {
		defer close(quittableChan)
		defer close(quit)
		for {
			select {
			case message := <-quittableChan:
				fmt.Println(message)
			case <-quit:
				fmt.Println("Quittable generator is being closed, worker goroutine rejoining main thread.")
				quit <- struct{}{}
				return
			}
		}
	}()
	return quittableChan
}

// This is to demonstrate that a generator can spawn multiple worker goroutines to handle requests.
// This example is easier to work with because the workers are producers for the channel.
func multiWorkerGenerator(message string, workerCount uint) <-chan string {
	generatorChan := make(chan string)
	for i := 0; uint(i) < workerCount; i++ {
		go func(i int) {
			for {
				generatorChan <- fmt.Sprintf("goroutine %v says %s", i + 1, message)
				time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			}
		}(i)
	}
	return generatorChan
}

func main() {
	boringChan := boringGenerator("It's dangerous to go alone! Take this.")
	for i := 0; i < 3; i++ {
		fmt.Println(<-boringChan)
	}
	fmt.Println()

	squareChan := squareGenerator()
	squareChan <- -3
	fmt.Println("squareGenerator processes the square of -3 to be:", <-squareChan)
	fmt.Println()

	quit := make(chan struct{})
	quittableChan := quittableGenerator(quit)
	for i := 0; i < 3; i++ {
		quittableChan <- "I quit!"
	}
	quit <- struct{}{}
	<-quit
	fmt.Println("Quittable worker should have just rejoined...")
	_, ok := <-quittableChan
	fmt.Println("quittableChan successfully closed:", !ok)
	fmt.Println()

	multiWorkerGeneratorChan := multiWorkerGenerator("It's dangerous to go alone! That's why we have 7 of us.", 7)
	for i := 0; i < 25; i++ {
		fmt.Println(<-multiWorkerGeneratorChan)
	}

}


// Additional improvements:
// 1) request / response pattern to support multiple clients and servers
//        we could either use separate channels for requests and responses with some reconciliation mechanism
//        or use a chan chan <Request> for communication (receive-only <-chan chan<- <Request> on the worker side)
// 2) quit channels to support multiple workers via fan-in fan-out pattern
