package main

import (
	micro "github.com/micro/go-micro"
	"github.com/peizhong/letsgo/pkg/log"
)

func main() {
	service := micro.NewService(micro.Name("greeter.client"))
	_ = service
	log.Info("aaa")
}
