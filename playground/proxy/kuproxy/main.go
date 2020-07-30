package main

import (
	"fmt"

	_ "net/http/pprof"

	_ "go.uber.org/automaxprocs"

	"github.com/peizhong/letsgo/playground/proxy/kuproxy/api"
	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
)

var (
	HostIP   = ""
	HttpPort = 8081
	TcpPort  = 8080
)

func main() {
	httpAddr := fmt.Sprintf("%s:%d", HostIP, HttpPort)
	tcpAddr := fmt.Sprintf("%s:%d", HostIP, TcpPort)

	rt := proxy.NewRuntime()
	go api.StartHttp(httpAddr, rt)
	if err := proxy.StartProxy(tcpAddr, rt); err != nil {
		panic(err)
	}
}
