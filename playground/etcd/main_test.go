package main

import (
	"github.com/peizhong/letsgo/pkg/log"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		<-time.After(time.Second)
		// 关闭一个已关闭的chan会panic

		// 如果没有goroutine引用chan，通道也会被垃圾回收，不一定要close
		close(ch)
	}()
	select {
	case <-time.After(time.Second * 2):
		log.Info("time")
	case <-ch:
		log.Info("ok")
	}
}
