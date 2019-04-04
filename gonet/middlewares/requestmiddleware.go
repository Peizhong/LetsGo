package middlewares

import "github.com/peizhong/letsgo/gonet"

type RequestMiddleware struct {
}

func (m *RequestMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	if next != nil {
		nch := make(chan struct{})
		go next(c, nch)
		<-nch
	}
	c.SayHi("requestMiddleware处理了上一个，然后再处理自己的")
	if ch != nil {
		close(ch)
	}
	return nil
}
