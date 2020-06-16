package main

import (
	"context"
	"time"
)

func UseContext() {

	bctx := context.Background()
	vcx := context.WithValue(bctx, "hi", "you")
	v := vcx.Value("hi")
	println(v)
	tx, tcan := context.WithTimeout(vcx, time.Minute)
	cc, can := context.WithCancel(tx)
	_ = cc
	_, _ = tcan, can
}
