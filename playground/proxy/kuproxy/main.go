package main

import (
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	_ "go.uber.org/automaxprocs"

	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HostIP   = ""
	HttpPort = 8081
	TcpPort  = 8080
)

func startHttp(httpAddr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.Handle("/admin", &adminHandler{})
	http.Handle("/metrics", promhttp.Handler())

	log.Println("http listen at", httpAddr)
	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	httpAddr := fmt.Sprintf("%s:%d", HostIP, HttpPort)
	tcpAddr := fmt.Sprintf("%s:%d", HostIP, TcpPort)

	go startHttp(httpAddr)

	if err := proxy.StartProxy(tcpAddr); err != nil {
		panic(err)
	}
}
