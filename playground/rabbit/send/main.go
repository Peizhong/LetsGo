package main

import (
	"bufio"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Sender interface {
	init(addr, queue string)
	publish(queue, text string) error
	call(queue, text string) (string, error)
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
	// 如果已经有了，可以不用再Declare的
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

func (r *rabbit) publish(queue, text string) error {
	err := r.channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(text),
		})
	return err
}

func (r *rabbit) call(queue, text string) (string, error) {
	err := r.channel.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: randomString(32),
			ReplyTo:       queue,
			Body:          []byte(text),
		})
	if err != nil {
		return "", err
	}
	// 还要开个consume接收queue的数据，用CorrelationId匹配
	return "", nil
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
		err := sender.publish("hello", text)
		if err != nil {
			log.Fatal(err)
		}
	}
	sender.close()
	log.Println("done")
}
