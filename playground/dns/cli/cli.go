package main

import (
	"fmt"
	"github.com/micro/mdns"
	"time"
)

func main() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 1)

loop:
	for {
		select {
		case entry := <-entriesCh:
			fmt.Printf("Got new entry: %v\n", entry)
			break loop
		default:
			// or _foobar._tcp
			if err := mdns.Lookup("topic:example.topic.pubsub.1", entriesCh); err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("found nothing")
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("bye")
}
