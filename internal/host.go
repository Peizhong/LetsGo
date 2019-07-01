package internal

import (
	"context"
	"log"
	"os"
	"os/signal"
)

// 强制退出/自动运行结束
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

func HostWithContext(app func(context.Context, chan struct{})) {
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan struct{}, 1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Println("start")
	go app(ctx, exit)
loop:
	select {
	case <-c:
		log.Println("force closing...")
		cancel()
		goto loop
	case <-exit:
		log.Println("bye")
	}
}