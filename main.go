package main

import (
	"os"
	"os/signal"

	"github.com/peizhong/letsgo/framework/log"
)

func main() {
	log.Info("Let's Go start")
	c := make(chan os.Signal, 1)
	go func() {

	}()
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Let's Go exit")
}
