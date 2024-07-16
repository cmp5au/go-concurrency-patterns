package main

import (
	"fmt"
)

type Topic struct{
	name string
	topicChan chan string
}

type Event struct{
	topic Topic
	message string
}

interface Publisher {
	PublishEvent(Event)
}

type Subscriber struct {
	GetSubscribedEvents() []Event
}

func main() {
	fmt.Println("Starting...")



	fmt.Println("Done!")
}
