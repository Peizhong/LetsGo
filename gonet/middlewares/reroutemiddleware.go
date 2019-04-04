package middlewares

import (
	"github.com/peizhong/letsgo/gonet"
)

type ReRouteMiddleware struct {
}

func (m *ReRouteMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	if next != nil {
		nch := make(chan struct{})
		go next(c, nch)
		<-nch
	}
	c.SayHi("ReRouteMiddleware要确定路径了哦")
	for _, route := range c.GetConfig().Routes {
		match := (*route).UpstreamRegexp.Match([]byte(c.SrcPath))
		if match {
			c.SayHi("ReRouteMiddleware配对成功")
			break
		}
	}
	if ch != nil {
		close(ch)
	}
	return nil
}
