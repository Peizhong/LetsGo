package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func sub() {
	redis := redis.NewClient(&redis.Options{
		Addr:     "193.112.41.28:6379",
		Password: "ur@hello123",
		DB:       0,
	})
	s := redis.Subscribe("mychannel1")

	// Wait for confirmation that subscription is created before publishing anything
	_, err := s.Receive()
	if err != nil {
		log.Panic(err)
	}

	// Go channel which receives messages.
	ch := s.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func main() {
	sub()
}
