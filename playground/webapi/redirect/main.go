package main

import (
	"context"
	"github.com/micro/go-micro"
	api "github.com/micro/micro/api/proto"
	"github.com/peizhong/letsgo/internal"
	"log"
)

type Redirect struct{}

func (r *Redirect) Url(ctx context.Context, req *api.Request, rsp *api.Response) error {
	rsp.StatusCode = int32(301)
	rsp.Header = map[string]*api.Pair{
		"Location": {
			Key:    "Location",
			Values: []string{"https://www.baidu.com"},
		},
	}
	return nil
}

func redirect(ctx context.Context, exit chan struct{}) {
	service := micro.NewService(
		micro.Name("go.micro.api.redirect"),
		micro.Context(ctx),
	)

	// parse command line flags
	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(new(Redirect)),
	)

	if err := service.Run(); err != nil {
		log.Println(err)
	}

	exit <- struct{}{}
}

func main() {
	internal.HostWithContext(redirect)
}

