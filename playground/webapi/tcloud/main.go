package main

import (
	"github.com/peizhong/letsgo/internal"
)

const (
	UseTcp = true

	IP = ""
	// 客户端访问的端口
	GatewayPort = 8080
	// 实际服务的端口
	RealPort = 8000
)

type redirect interface {
	start()
	stop()
}

func main() {
	var r redirect
	if UseTcp {
		r = &tcpRedirect{}
	} else {
		r = &httpRedirect{}
	}
	internal.Host(r.start, r.stop)
}
