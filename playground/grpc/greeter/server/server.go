package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/peizhong/letsgo/pkg/log"
	proto "github.com/peizhong/letsgo/playground/grpc/greeter"
	"github.com/peizhong/letsgo/playground/grpc/pubsub"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

// All methods of Sub will be executed when
// a message is received
type Sub struct{}

// Method can be of any name
func (s *Sub) Process(ctx context.Context, event *pubsub.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Infof("[pubsub.1] Received event %+v with metadata %+v\n", event, md)
	// do something with event
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
	// micro.Name("greeter"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// register subscriber
	micro.RegisterSubscriber("example.topic.pubsub.1", service.Server(), new(Sub))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
