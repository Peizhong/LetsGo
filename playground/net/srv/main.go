package main

import (
	"log"
	"net"
)

const (
	BUFFER_SIZE = 1024
)

func main() {
	ip := net.ParseIP("127.0.0.1")
	tcpAddr := &net.TCPAddr{IP: ip, Port: 0}
	// l, err := net.Listen("tcp", ":8081")
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	// port=0会自动分配
	log.Println(l.Addr().(*net.TCPAddr).Port)
	defer l.Close()
	for {
		if conn, err := l.Accept(); err == nil {
			gcon := conn
			go func() {
				log.Println("connect", gcon.RemoteAddr())
				for {
					var len int
					readbuf := make([]byte, BUFFER_SIZE)
					if _, err := gcon.Write([]byte("hi")); err != nil {
						log.Println(err)
						return
					}
					if len, err = gcon.Read(readbuf); err != nil {
						log.Println(err)
						return
					}
					log.Println(string(readbuf[:len]))
					if readbuf[0] == 0x01 {
						gcon.Close()
						return
					}
				}
			}()
		}
	}
	// 一个goroutine读
	// 一个goroutine写
}
