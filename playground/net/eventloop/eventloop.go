package main

import (
	"log"
	"net"
	"sync"
	"syscall"

	"github.com/panjf2000/gnet"
	"github.com/peizhong/letsgo/internal"
	"github.com/tidwall/evio"
)

const (
	BUFFER_SIZE = 1024
)

func nethttp() {
	l, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	f, err := l.(*net.TCPListener).File()
	if err != nil {
		panic(err)
	}
	fd := int(f.Fd())
	log.Println(fd)
	opened := func() {

	}
	readbuf := make([]byte, BUFFER_SIZE)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		opened()
		n, _ := conn.Read(readbuf)
		conn.Write(readbuf[:n])
		conn.Close()
	}
}

func miEpoll() {
	l, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	f, err := l.(*net.TCPListener).File()
	if err != nil {
		panic(err)
	}
	fd := int(f.Fd())
	syscall.SetNonblock(fd, true)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		idx := i
		// 每个协程创建一个epoll，监听同一个listener，不会冲突吗？
		p, err := syscall.EpollCreate1(0)
		if err != nil {
			panic(err)
		}
		if err := syscall.EpollCtl(p, syscall.EPOLL_CTL_ADD, int(fd),
			&syscall.EpollEvent{Fd: int32(fd),
				Events: syscall.EPOLLIN,
			},
		); err != nil {
			panic(err)
		}
		go func() {
			// wait
			events := make([]syscall.EpollEvent, 64)
			for {
				// 多个都会收到
				n, err := syscall.EpollWait(p, events, 1000)
				if err != nil && err != syscall.EINTR {
					panic(err)
				}
				// log.Println("idx", idx, "has", n)
				// 如果响应了多个，要按顺序执行
				for i := 0; i < n; i++ {
					cfd, _, err := syscall.Accept(int(events[i].Fd))
					if err != nil {
						// 多个epoll_wait会响应，有1个accept后，其他的accept返回EAGAIN
						if err == syscall.EAGAIN {
							log.Println("idx", idx, "lost", err.Error())
						} else {
							panic(err)
						}
					} else {
						log.Println("idx", idx, "got", cfd)
						// accpect后EpollCtl，只有自己的epoll_wait才会响应
						if err := syscall.EpollCtl(p, syscall.EPOLL_CTL_MOD, int(cfd),
							&syscall.EpollEvent{Fd: int32(cfd),
								Events: syscall.EPOLLIN | syscall.EPOLLOUT,
							},
						); err != nil {
							panic(err)
						}
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func goEvio() {
	var events evio.Events
	// Multithreaded: equal to runtime.NumProcs()
	events.NumLoops = -1
	// distribute the incoming connections between multiple loops
	// events.LoadBalance = evio.RoundRobin
	events.Serving = func(server evio.Server) (action evio.Action) {
		log.Println("serving", server.Addrs)
		return
	}
	// when a new connection has opened
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		// Data的数据不一定能接收完，每一个连接fd有一个context
		c.SetContext(struct{}{})
		// 系统层面设置，面向网络
		// opts.TCPKeepAlive = time.Second
		return
	}
	// when the server receives new data from a connection
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		// log.Println(string(in))
		out = in
		// action = evio.Detach
		return
	}
	if err := evio.Serve(events, "tcp://localhost:5000"); err != nil {
		panic(err.Error())
	}
}

type echoServer struct {
	*gnet.EventServer
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

func goNet() {
	echo := &echoServer{}
	gnet.Serve(echo, "tcp://localhost:5000", gnet.WithMulticore(true))
}

func main() {
	// GODEBUG=schedtrace=500 ./eventloop
	internal.PProf(goNet, nil)
}
