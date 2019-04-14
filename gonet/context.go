package gonet

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

type headers map[string][]string

type GatewayResponse struct {
	Headers headers
	Body    *[]byte
}

func NewGatewayResponse() *GatewayResponse {
	response := &GatewayResponse{
		Headers: make(headers),
	}
	return response
}

func (gw *GatewayResponse) AddHeader(key, value string) {
	if key == "" {
		return
	}
	if values, exist := gw.Headers[key]; exist {
		values = append(values, value)
	} else {
		gw.Headers[key] = []string{value}
	}
}

type Context struct {
	SrcPath    string `Remakr:"原始请求路径"`
	DstPath    string `Remark:"转发地址"`
	SrcHeaders headers

	Request   *http.Request
	ReqsCount uint64 `remark:"请求次数"`

	Responser http.ResponseWriter
	Response  *GatewayResponse

	Tracer opentracing.Tracer
}

func (c *Context) GetConfig() GatewayConfig {
	return gatewayConfig
}

func (c *Context) SayHi(message string) {
	// fmt.Println(time.Now(), message)
}
