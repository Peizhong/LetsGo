package main

import (
	"context"
	"fmt"
	"github.com/micro/mdns"
	"github.com/peizhong/letsgo/internal"
	"time"
)

func app(ctx context.Context, exit chan struct{}) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 1)
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("got cancel request, bye")
			break loop
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
	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(app)
}
