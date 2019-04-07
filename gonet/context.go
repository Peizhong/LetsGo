package gonet

import (
	"net/http"
)

type headers map[string]string

type Context struct {
	SrcPath    string `Remakr:"原始请求路径"`
	DstPath    string `Remark:"转发地址"`
	SrcHeaders headers
	DstHeaders headers

	ReqsCount uint64 `remark:"请求次数"`

	Responser http.ResponseWriter
	Response  *[]byte
}

func (c *Context) GetConfig() GatewayConfig {
	return gatewayConfig
}

func (c *Context) SayHi(message string) {
	// fmt.Println(time.Now(), message)
}
