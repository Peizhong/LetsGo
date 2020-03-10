package main

import (
	"fmt"
	"github.com/peizhong/letsgo/internal"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type tcpRedirect struct {
	listener *net.TCPListener
	connMap  map[string]*connMatch
	m        sync.RWMutex
}

var (
	gatewayAddr, _ = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", IP, GatewayPort))
	realAddr, _    = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", IP, RealPort))
)

type connMatch struct {
	outside       *net.TCPConn //8007 tcp链路 accept
	inside        *net.TCPConn //8008 tcp链路 tunnel
	acceptAddTime int64        //接受请求的时间
	count         int
}

func (t *tcpRedirect) start() {
	joinConn := func(conn1 *net.TCPConn, conn2 *net.TCPConn) {
		f := func(local *net.TCPConn, remote *net.TCPConn) {
			//defer保证close
			defer local.Close()
			defer remote.Close()
			//使用io.Copy传输两个tcp连接
			io.Copy(local, remote)
			// internal.CheckError(err, "io.Copy")
			log.Println("conn end", local.LocalAddr(), local.RemoteAddr(), remote.LocalAddr(), remote.RemoteAddr())
		}
		go f(conn2, conn1)
		go f(conn1, conn2)
	}
	addConnMathAccept := func(accept *net.TCPConn) {
		recordGeo(accept.RemoteAddr().String())
		key := getIP(accept.RemoteAddr().String())
		t.m.Lock()
		if v, ok := t.connMap[key]; ok {
			v.count++
			v.inside.Close()
		}
		t.m.Unlock()
		inside, err := net.DialTCP("tcp", nil, realAddr)
		if internal.CheckError(err, "DialTCP") != nil {
			return
		}
		match := &connMatch{
			outside:       accept,
			inside:        inside,
			acceptAddTime: time.Now().Unix(),
		}
		t.m.Lock()
		t.connMap[key] = match
		t.m.Unlock()
		joinConn(match.inside, match.outside)
	}
	err := prepareGeo()
	internal.CheckError(err, "prepareGeo")
	// 开启监听
	t.connMap = make(map[string]*connMatch)
	t.listener, err = net.ListenTCP("tcp", gatewayAddr)
	internal.CheckError(err, "ListenTCP")
	for {
		// 与客户端建立连接，然后realport建立连接
		tcpConn, err := t.listener.AcceptTCP()
		if internal.CheckError(err, "AcceptTCP") != nil {
			continue
		}
		// 交换双方的数据
		addConnMathAccept(tcpConn)
	}
}

func (t *tcpRedirect) stop() {
	log.Println("exit geo")
	err := unloadGeo()
	internal.CheckError(err, "Unload Geo")
}
