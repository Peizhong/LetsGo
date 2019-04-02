package gonet

import (
	"fmt"
)

type headers map[string]string

type Context struct {
	srcPath    string `Remakr:"原始请求路径"`
	dstPath    string `Remark:"转发地址"`
	srcHeaders headers
	dstHeaders headers
}

func (*Context) SayHi(message string) {
	fmt.Println(message)
}
