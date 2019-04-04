package middlewares

import (
	"fmt"

	"github.com/peizhong/letsgo/gonet"
)

type ResponseMiddleware struct {
}

func (m *ResponseMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	if next != nil {
		nch := make(chan struct{})
		go next(c, nch)
		<-nch
	}
	c.SayHi("responseMiddleware处理了上一个，然后再处理自己的")
	fmt.Fprintf(c.Responser, "Hello,"+"这是最后一个咯"+c.SrcPath)
	if ch != nil {
		close(ch)
	}
	return nil
}
