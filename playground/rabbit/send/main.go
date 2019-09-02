package main

import (
	"bufio"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func send() {
	conn, err := amqp.Dial("amqp://guest:guest@193.112.41.28:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("online on, write something")
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		if text == "exit" {
			break
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body:        []byte(text),
			})
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("done")
}

func main() {
	send()
}
