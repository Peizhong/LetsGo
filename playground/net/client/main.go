package main

import (
	"github.com/peizhong/letsgo/internal"
	"log"
	"net"
)

const (
	Port = 8080
)

func main() {
	internal.Host(func() {
		ip := net.ParseIP("127.0.0.1")
		tcpAddr := &net.TCPAddr{IP: ip, Port: Port}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		for {
			_, err = conn.Write([]byte("hello"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}, nil)
}
