package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/peizhong/letsgo/pkg/log"
	"time"

	"context"
	"github.com/micro/go-micro"
	proto "github.com/peizhong/letsgo/playground/grpc/pubsub"
)

// send events using the publisher
func sendEv(topic string, p micro.Publisher) {
	t := time.NewTicker(time.Second)

	for _ = range t.C {
		// create new event
		ev := &proto.Event{
			Id:        uuid.New().String(),
			Timestamp: time.Now().Unix(),
			Message:   fmt.Sprintf("Messaging you all day on %s", topic),
		}

		log.Infof("publishing %+v\n", ev)

		// publish an event
		if err := p.Publish(context.Background(), ev); err != nil {
			log.Infof("error publishing: %v", err)
		}
		t.Stop()
	}
}

func main() {
	// create a service
	service := micro.NewService(
		micro.Name("go.micro.cli.pubsub"),
	)
	// parse command line
	service.Init()

	// create publisher
	pub1 := micro.NewPublisher("example.topic.pubsub.1", service.Client())
	// pub2 := micro.NewPublisher("example.topic.pubsub.2", service.Client())

	// pub to topic 1
	go sendEv("example.topic.pubsub.1", pub1)
	// pub to topic 2
	// go sendEv("example.topic.pubsub.2", pub2)

	// block forever
	select {}
}
