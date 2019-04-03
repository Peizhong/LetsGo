package middlewares

import "github.com/peizhong/letsgo/gonet"

type RequestMiddleware struct {
}

func (m *RequestMiddleware) Invoke(c *gonet.Context, next func(*gonet.Context) error) (err error) {
	if next != nil {
		err = next(c)
	}
	c.SayHi("requestMiddleware处理了上一个，然后再处理自己的")
	return
}
