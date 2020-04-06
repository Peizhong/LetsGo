package main

import (
	"context"
	"time"
)

func UseContext() {
	ch := make(chan struct{})
	time.AfterFunc(time.Second*2, func() {
		ch <- struct{}{}
	})
	select {
	case <-ch:
		println("oo")
	}
	select {
	case <-ch:
		println("aa")
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	// parent的done完成后，cc也会done
	cc, cancel := context.WithCancel(ctx)
	defer cancel()
	select {
	case <-ctx.Done():
		println("out")
	case <-cc.Done():
		println("ok")
	}
	select {
	case <-ctx.Done():
		println("out")
	case <-cc.Done():
		println("ok")
	}
	select {
	case <-ctx.Done():
		println("out")
	case <-cc.Done():
		println("ok")
	default:
		println("noo")
	}
}
