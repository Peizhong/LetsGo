package main

import (
	"os"
	"os/signal"

	"github.com/peizhong/letsgo/framework/log"
)

func main() {
	log.Info("Let's Go start")
	workers := loadWorkers()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Closing Services")
	stopWorkers(workers...)
	log.Info("Let's Go exit")
}
