package proxy

import (
	"log"
	"net"
)

func StartProxy(tcpAddr string) error {
	log.Println("tcp listen at", tcpAddr)
	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		return err
	}
	go cleanUpBrokenConn()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		if err = makeTunnel(conn); err != nil {
			log.Panicln("abort tcp conn", err.Error())
			conn.Close()
		}
	}
}
