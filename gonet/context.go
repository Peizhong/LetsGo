package gonet

import (
	"fmt"
	"net/http"
	"time"
)

type headers map[string]string

type Context struct {
	SrcPath    string `Remakr:"原始请求路径"`
	DstPath    string `Remark:"转发地址"`
	SrcHeaders headers
	DstHeaders headers

	Responser http.ResponseWriter
}

func (c *Context) SayHi(message string) {
	fmt.Println(time.Now(), message)
}
