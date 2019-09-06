package main

import (
	"bufio"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Sender interface {
	init(addr, queue string)
	publish(text string) error
	close()
}

type rabbit struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func (r *rabbit) init(addr, queue string) {
	var err error
	r.conn, err = amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	r.queue, err = r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *rabbit) publish(text string) error {
	err := r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(text),
		})
	return err
}

func (r *rabbit) close() {
	r.channel.Close()
	r.conn.Close()
}

func main() {
	var sender Sender
	sender = &rabbit{}
	sender.init("amqp://guest:guest@193.112.41.28:5672/", "hello")
	log.Println("online on, write something")
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		if text == "exit" {
			break
		}
		err := sender.publish(text)
		if err != nil {
			log.Fatal(err)
		}
	}
	sender.close()
	log.Println("done")
}
