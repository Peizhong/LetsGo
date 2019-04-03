package gonet

import (
	"sync"
)

type IMiddleware interface {
	Invoke(*Context, func(*Context) error) error
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
	var next func(*Context) error
	for _, md := range middlewares {
		backUpNext := next
		bakcupMd := md
		nextNext := func(c *Context) error {
			err := bakcupMd.Invoke(c, backUpNext)
			return err
		}
		next = nextNext
	}
	entry = next
	return
}
