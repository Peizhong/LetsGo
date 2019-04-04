package middlewares

import "github.com/peizhong/letsgo/gonet"

type FirstMiddleware struct {
}

func (m *FirstMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	c.SayHi("最后一个中间件，没有next")
	if ch != nil {
		close(ch)
	}
	return nil
}
