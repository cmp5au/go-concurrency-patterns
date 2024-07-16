package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Response struct {
	id int
}

func main() {
	fmt.Println("Starting...")

	ctx := context.Background()

	id := 10

	for i := 0; i < 10; i++ {
		fmt.Println(thirdPartyCallWithTimeoutWrapper(ctx, id))
	}

	fmt.Println("Done!")
}

func thirdPartyCallWithTimeoutWrapper(ctx context.Context, id int) string {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
	defer cancel()

	ch := make(chan Response)

	go expensiveThirdPartyAPICall(id, ch)

	for {
		select {
		case <-ctxWithTimeout.Done():
			return "Timed out"
		case val := <-ch:
			return fmt.Sprintf("Got %v", val)
		}
	}
}


func expensiveThirdPartyAPICall(id int, ch chan Response) {
	time.Sleep(time.Duration(195 + rand.Intn(10)) * time.Millisecond)

	ch <- Response{id}
}
