package main

import (
	"context"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/peizhong/letsgo/internal"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func subscribe(addr, queue string) (msgs <-chan amqp.Delivery, closer func()) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	// open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	msgs, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return msgs, func() {
		ch.Close()
		conn.Close()
	}
}

func recv(ctx context.Context, exit chan struct{}) {
	msgs, closer := subscribe("amqp://guest:guest@193.112.41.28:5672/", "hello")
	buf := linkedliststack.New()
	t := time.NewTimer(5 * time.Second)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case d := <-msgs:
			buf.Push(d)
		case <-t.C:
			log.Println("clear", buf.Size())
			buf.Clear()
			t.Reset(5 * time.Second)
		}
	}
	closer()
	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(recv)
}
