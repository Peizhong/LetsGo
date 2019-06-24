package main

import (
	"context"
	"fmt"
	"github.com/peizhong/letsgo/internal"

	"github.com/micro/go-micro"
	proto "github.com/peizhong/letsgo/playground/grpc/greeter"
)

func client(ctx context.Context, exit chan struct{}) {
	// create service
	service := micro.NewService(
		micro.Name("greeter.client"),
		// with our cancellation context
		micro.Context(ctx),
	)
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)

	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(client)
}
