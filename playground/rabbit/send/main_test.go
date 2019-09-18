package main

import (
	"fmt"
	"log"
	"testing"
)

func BenchmarkDoRabbit(b *testing.B) {
	var sender Sender
	sender = &rabbit{}
	sender.init("amqp://guest:guest@193.112.41.28:5672/", "hello")
	for i := 0; i < b.N; i++ {
		text := fmt.Sprintf("message %d", i)
		err := sender.publish(text,"hello")
		if err != nil {
			log.Fatal(err)
		}
	}
	sender.close()
}
