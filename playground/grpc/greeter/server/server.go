package main

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/micro/go-micro"
	"github.com/peizhong/letsgo/pkg/mock_foo"
	proto "github.com/peizhong/letsgo/playground/grpc/greeter"
	"log"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// mock
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mreg := mock_foo.NewMockRegistry(ctrl)
	//String()
	mreg.EXPECT().
		String().
		DoAndReturn(func() string {
			log.Println("Mock Registry.String")
			return "Mock"
		})

	mreg.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		DoAndReturn(func(args ...interface{}) error {
			log.Println("Mock Registry.Register")
			return nil
		}).
		AnyTimes()

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
