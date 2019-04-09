package middlewares

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/peizhong/letsgo/gonet"
)

type FirstMiddleware struct {
}

func (m *FirstMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	c.SayHi("最后一个中间件，没有next")
	set := hashset.New() // empty
	set.Add(1)           // 1
	set.Clear()          // empty
	if ch != nil {
		close(ch)
	}
	return nil
}
