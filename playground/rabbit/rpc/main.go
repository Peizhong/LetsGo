package main

import (
	"context"
	"github.com/peizhong/letsgo/internal"
	"github.com/streadway/amqp"
	"log"
)

func subscribe(addr, queue string) (ch *amqp.Channel, msgs <-chan amqp.Delivery, closer func()) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	// open a channel
	ch, err = conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack, call Delivery.Ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return ch, msgs, func() {
		ch.Close()
		conn.Close()
	}
}

func recv(ctx context.Context, exit chan struct{}) {
	ch, msgs, closer := subscribe("amqp://guest:guest@193.112.41.28:5672/", "hello")
loop:
	for {
		select {
		case d := <-msgs:
			err := ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("gotcha"),
				})
			if err == nil {
				d.Ack(false)
			}
		case <-ctx.Done():
			break loop
		}
	}
	closer()
	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(recv)
}
