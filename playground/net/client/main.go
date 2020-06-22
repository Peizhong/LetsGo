package main

import (
	"log"
	"net"
	"net/http"

	"github.com/peizhong/letsgo/internal"
)

const (
	Port = 8080
)

func httpclient() {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}
	resp, err = client.Get("http://golang.org")
}

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
