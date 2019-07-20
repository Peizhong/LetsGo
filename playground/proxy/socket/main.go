package main

import (
	"github.com/peizhong/letsgo/internal"
	"io"
	"log"
	"net"
	"os"
)

func dail() {
	addr := "www.baidu.com:80"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Printf("local: %s, remote: %s", conn.LocalAddr().String(), conn.RemoteAddr().String())
	n, err := conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("write length: %d", n)
	io.Copy(os.Stdout, conn)
}

func start()  {
	dail()
}



func main() {
	internal.Host(start, nil)
}
