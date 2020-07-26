package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	_ "net/http/pprof"

	_ "go.uber.org/automaxprocs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HostIP   = ""
	HttpPort = 8081
	TcpPort  = 8080
)

func startHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	http.Handle("/admin", &adminHandler{})
	http.Handle("/metrics", promhttp.Handler())

	httpAddr := fmt.Sprintf("%s:%d", HostIP, HttpPort)
	log.Println("http listen at", httpAddr)
	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	tcpAddr := fmt.Sprintf("%s:%d", HostIP, TcpPort)
	log.Println("tcp listen at", tcpAddr)
	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	go startHttp()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		bindUpstream(conn)
	}

}
