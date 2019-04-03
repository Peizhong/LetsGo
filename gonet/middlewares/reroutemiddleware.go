package middlewares

import "github.com/peizhong/letsgo/gonet"

type ReRouteMiddleware struct {
}

func (m *ReRouteMiddleware) Invoke(c *gonet.Context, next func(*gonet.Context) error) (err error) {
	if next != nil {
		err = next(c)
	}
	c.SayHi("ReRouteMiddleware要确定路径了哦")
	return
}
