package middlewares

import (
	"fmt"

	"github.com/peizhong/letsgo/gonet"
)

type ResponseMiddleware struct {
}

func (m *ResponseMiddleware) Invoke(c *gonet.Context, next func(*gonet.Context) error) (err error) {
	if next != nil {
		err = next(c)
	}
	c.SayHi("responseMiddleware处理了上一个，然后再处理自己的")
	fmt.Fprintf(c.Responser, "Hello,"+"这是最后一个咯"+c.SrcPath)
	return
}
