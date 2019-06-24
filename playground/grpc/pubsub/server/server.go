package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/server"
	"github.com/peizhong/letsgo/pkg/log"
	proto "github.com/peizhong/letsgo/playground/grpc/pubsub"

	"context"
)

// All methods of Sub will be executed when
// a message is received
type Sub struct{}

// Method can be of any name
func (s *Sub) Process(ctx context.Context, event *proto.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Infof("[pubsub.1] Received event %+v with metadata %+v\n", event, md)
	// do something with event
	return nil
}

// Alternatively a function can be used
func subEv(ctx context.Context, event *proto.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Infof("[pubsub.2] Received event %+v with metadata %+v\n", event, md)
	// do something with event
	return nil
}

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"http://localhost:8500"}
	})

	// create a service
	service := micro.NewService(
		micro.Name("go.micro.srv.pubsub"),
		micro.Registry(reg),
	)

	// parse command line
	service.Init()

	// register subscriber
	micro.RegisterSubscriber("example.topic.pubsub.1", service.Server(), new(Sub))

	// register subscriber with queue, each message is delivered to a unique subscriber
	micro.RegisterSubscriber("example.topic.pubsub.2", service.Server(), subEv, server.SubscriberQueue("queue.pubsub"))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
