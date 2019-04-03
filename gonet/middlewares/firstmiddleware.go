package middlewares

import "github.com/peizhong/letsgo/gonet"

type FirstMiddleware struct {
}

func (m *FirstMiddleware) Invoke(c *gonet.Context, next func(*gonet.Context) error) (err error) {
	c.SayHi("最后一个中间件，没有next")
	return
}
