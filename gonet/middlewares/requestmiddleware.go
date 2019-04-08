package middlewares

import (
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/peizhong/letsgo/gonet"
)

type RequestMiddleware struct {
}

var headerDontSend = map[string]interface{}{
	"Set-Cookie": 1,
	"2":          2,
	"3":          3,
}

var headerDontRecv = map[string]interface{}{
	"Set-Cookie": 1,
	"2":          2,
	"3":          3,
}

func (m *RequestMiddleware) Invoke(c *gonet.Context, ch chan<- struct{}, next func(*gonet.Context, chan<- struct{}) error) (err error) {
	if next != nil {
		nch := make(chan struct{})
		go next(c, nch)
		<-nch
	}
	c.SayHi("requestMiddleware处理了上一个，然后再处理自己的")
	atomic.AddUint64(&c.ReqsCount, 1)
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://www.baidu.com/s?wd="+c.SrcPath, nil)
	if err == nil {
		response, err := client.Do(request)
		if err == nil {
			for k, values := range response.Header {
				if _, exist := headerDontRecv[k]; exist {
					// continue
				}
				for _, v := range values {
					c.Response.AddHeader(k, v)
				}
			}
			bytes, err := ioutil.ReadAll(response.Body)
			if err == nil {
				c.Response.AddHeader("Content-Length", string(len(bytes)))
				c.Response.Body = &bytes
			}
		}
	}
	if ch != nil {
		close(ch)
	}
	return err
}
