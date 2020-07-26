package main

import (
	"io"
	"log"
	"net"
	"strings"
)

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	ok = true
	method = line[:s1]
	s2 += s1 + 1
	requestURI = line[s1+1 : s2]
	proto = line[s2+1:]
	return
}

func bindUpstream(down net.Conn) {
	_selector, _ := getSelector()
	buf := make([]byte, 1024)
	n, err := down.Read(buf)
	if err != nil {
		return
	}
	// 检查是否为http
	_, url, _, ok := parseRequestLine(string(buf))
	if !ok {
		return
	}
	log.Println("new conn", string(buf))
	upconn, err := net.Dial("tcp", _selector.TellMe(url))
	if err != nil {
		log.Println("upconn err", err.Error())
		return
	}
	// 将截取的http写到upstream
	// wrap = bufio.NewReader()
	_, err = upconn.Write(buf[:n])
	if err != nil {
		return
	}
	// A -> B
	go func() {
		defer down.Close()
		defer upconn.Close()
		_, err := io.Copy(down, upconn)
		if err != nil {
			return
		}
	}()
	// B -> A
	go func() {
		defer down.Close()
		defer upconn.Close()
		_, err = io.Copy(upconn, down)
		if err != nil {
			return
		}
	}()
}
