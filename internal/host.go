package internal

import (
	"log"
	"os"
	"os/signal"
)

func Host(start, stop func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Println("start")
	go start()
	<-c
	log.Println("closing...")
	stop()
	log.Println("bye")
}
