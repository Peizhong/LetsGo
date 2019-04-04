package gonet

import (
	"sync"
)

type IMiddleware interface {
	// chan struct{} 不能输入数据，通过close关闭
	Invoke(*Context, chan<- struct{}, func(*Context, chan<- struct{}) error) error
}

var m sync.Mutex

var middlewares []IMiddleware

func AddMiddleware(middleware IMiddleware) {
	m.Lock()
	middlewares = append(middlewares, middleware)
	m.Unlock()
}

func BuildPipeline() (entry func(*Context) error) {
	// build
	var next func(*Context, chan<- struct{}) error
	for _, md := range middlewares {
		backUpNext := next
		bakcupMd := md
		nextNext := func(c *Context, ch chan<- struct{}) error {
			err := bakcupMd.Invoke(c, ch, backUpNext)
			return err
		}
		next = nextNext
	}
	entry = func(c *Context) error {
		return next(c, nil)
	}
	return
}
