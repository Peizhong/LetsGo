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
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	  )
	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	msgs, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack, call Delivery.Ack
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
			d.Ack(false)
			// d.Nack(false,true) 进死信
			// d.Reject(false)
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
