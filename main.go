package main

import (
	"letsgo/framework/config"
	"os"
	"os/signal"

	"letsgo/framework/log"
)

func main() {
	log.Info("Let's Go start at", config.RunMode)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Closing Services")
	log.Info("Let's Go exit")
}
