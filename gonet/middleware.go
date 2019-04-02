package gonet

type MiddleWare interface {
	Invoke(*Context, func(*Context))
}

type MiddleWareBuilder struct {
	middlewares []MiddleWare
}

func (m *MiddleWareBuilder) Add(middleware MiddleWare) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *MiddleWareBuilder) Build() (entry func(*Context)) {
	// demo
	m.Add(&firstMiddleware{})
	m.Add(&requestMiddleware{})
	m.Add(&responseMiddleware{})

	// build
	var next func(*Context)
	for _, md := range m.middlewares {
		backUpNext := next
		bakcupMd := md
		nextNext := func(c *Context) {
			bakcupMd.Invoke(c, backUpNext)
		}
		next = nextNext
	}
	entry = next
	return
}

type firstMiddleware struct {
}

func (m *firstMiddleware) Invoke(c *Context, next func(*Context)) {
	c.SayHi("最后一个中间件，没有next")
}

type requestMiddleware struct {
}

func (m *requestMiddleware) Invoke(c *Context, next func(*Context)) {
	if next != nil {
		next(c)
	}
	c.SayHi("requestMiddleware处理了上一个，然后再处理自己的")
}

type responseMiddleware struct {
}

func (m *responseMiddleware) Invoke(c *Context, next func(*Context)) {
	if next != nil {
		next(c)
	}
	c.SayHi("responseMiddleware处理了上一个，然后再处理自己的")
}
