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
		f := func(tag string, local *net.TCPConn, remote *net.TCPConn) {
			//defer保证close
			defer local.Close()
			defer remote.Close()
			log.Println("establish", tag)
			//使用io.Copy传输两个tcp连接
			io.Copy(local, remote)
			// internal.CheckError(err, "io.Copy")
			log.Println("conn end", local.LocalAddr(), local.RemoteAddr(), remote.LocalAddr(), remote.RemoteAddr())
		}
		go f("out->in", conn2, conn1)
		go f("in->out", conn1, conn2)
	}
	addConnMathAccept := func(accept *net.TCPConn) {
		key := getIP(accept.RemoteAddr().String())
		recordGeo(key)
		var limit bool
		var cnt int
		t.m.Lock()
		if v, ok := t.connMap[key]; ok {
			cnt = v.count + 1
			if cnt > 10 {
				limit = true
			}
			v.outside.Close()
			v.inside.Close()
		}
		t.m.Unlock()
		if limit {
			log.Println("byb bye", accept.RemoteAddr().String())
			accept.Close()
			return
		}
		inside, err := net.DialTCP("tcp", nil, realAddr)
		if internal.CheckError(err, "DialTCP") != nil {
			return
		}
		match := &connMatch{
			outside:       accept,
			inside:        inside,
			acceptAddTime: time.Now().Unix(),
			count:         cnt,
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
	if internal.CheckError(err, "ListenTCP"); err != nil {
		return
	}
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
