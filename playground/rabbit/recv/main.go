package main

import (
	"context"
	"github.com/peizhong/letsgo/internal"
	"log"
	"github.com/streadway/amqp"
)

func recv(ctx context.Context,exit chan struct{}) {
	conn, err := amqp.Dial("amqp://guest:guest@193.112.41.28:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}
loop:
	for{
		select {
		case <-ctx.Done():
			break loop
		case d := <-msgs:
			log.Printf("Received a message: %s", d.Body)
			break
		}
	}
	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(recv)
}
