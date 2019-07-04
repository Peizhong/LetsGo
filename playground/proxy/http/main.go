package main

import (
	"bytes"
	"fmt"
	"github.com/peizhong/letsgo/internal"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

const ()

func start() {
	handle := func(client net.Conn) {
		buf := make([]byte, 1024)
		// 填充到buf, 要有长度
		n, err := client.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("request %d bytes:\r\n%s", n, string(buf))
		var method, host, address string
		fmt.Sscanf(string(buf[:bytes.IndexByte(buf, '\n')]), "%s%s", &method, &host)
		hostPortURL, err := url.Parse(host)
		if err != nil {
			log.Println(err)
			return
		}
		if hostPortURL.Opaque == "443" { //https访问
			address = hostPortURL.Scheme + ":443"
		} else { //http访问
			if strings.Index(hostPortURL.Host, ":") == -1 {
				//host不带端口， 默认80
				address = hostPortURL.Host + ":80"
			} else {
				address = hostPortURL.Host
			}
		}
		log.Println(method, host, address, hostPortURL)

		// 固定的服务地址
		const targetAddress = "localhost:8000"
		service, err := net.Dial("tcp", targetAddress)
		if err != nil {
			log.Panic(err)
		}
		// 如果时connect方法，要告诉客户端已经连接好
		if method == "CONNECT" {
			fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n")
		} else {
			n, err = service.Write(buf[:n])
		}
		var wg sync.WaitGroup
		wg.Add(2)
		// 流、块的区别 https://segmentfault.com/q/1010000009691632/
		// 流量转发 io.copy
		// https://studygolang.com/articles/8647
		go func() {
			w, _ := io.Copy(service, client)
			log.Println("client to service: ", w)
			wg.Done()
		}()
		go func() {
			header := make([]byte, 16)
			_, err := io.ReadAtLeast(service, header, 3)
			if err != nil {
				log.Panic(err)
			}
			log.Println("header: ", string(header))
			client.Write(header)
			w, _ := io.Copy(client, service)
			log.Println("service to client: ", w)
			wg.Done()
		}()
		wg.Wait()
		log.Println("proxy completed")
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		handle(client)
	}
}

func main() {
	type s struct {
	}
	internal.Host(start, nil)
}
