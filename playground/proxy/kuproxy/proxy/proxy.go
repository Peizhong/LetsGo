package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"syscall"
)

type proxy struct {
	rt *Runtime

	ln *net.TCPListener // all the listeners
}

const (
	RoomIdParam      = "roomId"
	ServiceNameParam = "serviceName"
)

type proxyEndMsg struct {
	service, endpoint, id string
}

var (
	ParseRequestError = errors.New("Parse request error")
	// 客户端连接断开后，通知清理
	brokenConnCh = make(chan proxyEndMsg, 10)
)

func StartProxy(tcpAddr string, rt *Runtime) error {
	log.Println("tcp listen at", tcpAddr)
	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		return err
	}
	tcpl := l.(*net.TCPListener)
	p := &proxy{
		rt: rt,
		ln: tcpl,
	}
	go p.cleanUpBrokenConn()
	p.stdserve()
	return nil
}

// 使用net标准库
func (p *proxy) stdserve() {
	for {
		conn, err := p.ln.Accept()
		if ce := checkError(err, "accept"); ce != nil {
			log.Println(ce.Error())
			continue
		}
		err = p.makeTunnel(conn)
		if ce := checkError(err, "make tunnel"); ce != nil {
			log.Println(ce.Error())
			conn.Close()
		}
	}
}

type epoll struct {
	fd int // epoll fd
}

// AddRead ...
func (p *epoll) AddRead(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

func (p *epoll) AddReadWrite(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

func (p *epoll) Wait(fn func(fd int) error) error {
	events := make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.fd, events, 100)
		if err != nil && err != syscall.EINTR {
			log.Println("wait err", err.Error())
			return err
		}
		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if err := fn(fd); err != nil {
				log.Println("fn err", err.Error())
				return err
			}
		}
	}
}

// 使用epoll
func (p *proxy) serve() {
	// 获得tcp监听的文件描述符
	file, err := p.ln.File()
	if err != nil {
		panic(err)
	}
	lnFd := int(file.Fd())
	workerCnt := runtime.GOMAXPROCS(0)
	var index = 0
	acceptBalancer := func(i int) bool {
		if i == (index % workerCnt) {
			index += 1
			return true
		}
		return false
	}
	log.Println("start ", workerCnt, "workers")
	for i := 0; i < 1; i++ {
		index := i
		go func() {
			// 设置epoll
			epollFd, err := syscall.EpollCreate1(0)
			if err != nil {
				panic(err)
			}
			ep := epoll{fd: epollFd}
			// 将tcp监听描述符设为非阻塞
			syscall.SetNonblock(lnFd, true)
			// 将tcp监听描述符添加到epoll
			ep.AddRead(lnFd)

			fdconns := make(map[int]int)
			ep.Wait(func(fd int) error {
				if !acceptBalancer(index) {
					// return nil
				}
				if _, ok := fdconns[fd]; !ok {
					// accept新的连接
					nfd, _, err := syscall.Accept(fd)
					if err != nil {
						log.Println("accept err", err)
						if err == syscall.EAGAIN {
							return nil
						}
						// 惊群，accept失败
						return nil
					}
					log.Println("accept", nfd)
					header := make([]byte, 512)
					n, err := syscall.Read(nfd, header)
					if err != nil {
						log.Println(err.Error())
					}
					if n < 0 {
						log.Println("can't read")
						return nil
					}
					header = header[:n]
					if err := syscall.SetNonblock(nfd, true); err != nil {
						log.Println("set noblocking err", err)
						return err
					}
					ep.AddRead(nfd)

					// do upstream
					upFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
					if err != nil {
						log.Println("socket err", err.Error())
						return err
					}
					err = syscall.Connect(upFd, &syscall.SockaddrInet4{
						Port: 3000,
						Addr: [...]byte{192, 168, 3, 143},
					})
					if err != nil {
						log.Println("conn err", err.Error())
						return err
					}
					log.Println("connect ", upFd)
					log.Println("write back", string(header))
					syscall.Write(upFd, header)
					if err := syscall.SetNonblock(upFd, true); err != nil {
						log.Println("set noblocking err", err)
						return err
					}
					fdconns[nfd] = upFd
					fdconns[upFd] = nfd
					ep.AddRead(upFd)
				} else {
					// 1. 直接发送  2. 看返回值, 没发送完才挂上EPOLLOUT  3. 发送完就把EPOLLOUT 取消掉
					in := make([]byte, 512)
					n, err := syscall.Read(fd, in)
					if n == 0 || err != nil {
						if err != nil {
							if err == syscall.EAGAIN {
							} else {
								log.Println("unknow read err", err.Error())
							}
						}
						//close(fd)后，epoll中注册的相关sock_fd，就会自动被清理？
						log.Println("close", fd)
						syscall.Close(fd)
						if tfd, ok := fdconns[fd]; ok {
							syscall.Close(tfd)
							delete(fdconns, tfd)
							log.Println("close2 ", tfd)
						}
						delete(fdconns, fd)
						return nil
					}
					in = in[:n]
					if othFd, ok := fdconns[fd]; ok {
						_, err := syscall.Write(othFd, in)
						if err != nil {
							log.Println("write err", err.Error())
							return err
						}
					}
					return nil
				}
				return nil
			})
		}()
	}
	select {}
}

// makeTunnel: 建立双向转发通道
// todo: 每个连接会产生2个goroutine，换成select/epoll?
func (p *proxy) makeTunnel(downConn net.Conn) error {
	buf := make([]byte, 512)
	n, err := downConn.Read(buf)
	if ce := checkError(err, "read downstream"); ce != nil {
		log.Println(ce.Error())
		return ce
	}
	// 检查请求是否是包含roomId的http请求
	_, url, roomId, serviceName := parseRequestLine(string(buf))
	if url == "" {
		log.Println("unknow request")
		return ParseRequestError
	}
	if roomId == "" {
		log.Println("no roomId, use default")
	}
	if serviceName == "" {
		// 调试用，应该在http请求中包含
		serviceName = DefaultServiceName
	}
	log.Println(url, roomId, serviceName)
	selector, ok := p.rt.SelectorMap[serviceName]
	if !ok {
		selector = NewSelector(serviceName, p.rt)
		p.rt.SelectorMap[serviceName] = selector
	}
	endpoint, newSelect, err := selector.SelectEndpoint(roomId)
	if ce := checkError(err, fmt.Sprintf("select %s endpoint", selector.ServiceName())); ce != nil {
		log.Println(ce.Error())
		return ce
	}
	log.Println("dial to ", endpoint)
	upConn, err := net.Dial("tcp", endpoint)
	if ce := checkError(err, fmt.Sprintf("dial %s", endpoint)); ce != nil {
		log.Println(ce.Error())
		// dial upstream失败，放弃本次连接
		return ce
	}
	// 将截取的tcp写回upstream
	_, err = upConn.Write(buf[:n])
	if ce := checkError(err, "write upstream"); ce != nil {
		log.Println(ce.Error())
		// 写入upstream失败，放弃本次连接
		upConn.Close()
		return ce
	}
	// 向公共缓存写入roomId和endpoint的信息
	err = selector.ConfirmEnpoint(endpoint, roomId, newSelect)
	if ce := checkError(err, "confirm endpoint"); ce != nil {
		// 如果写入失败，不能保证后续该roomid的连接到同一个endpoint
		log.Println("WARN", ce.Error())
		// 但仍继续往下走
		// return ce
	}
	// upstream -> downstream
	go func() {
		defer downConn.Close()
		defer upConn.Close()
		io.Copy(downConn, upConn)
	}()
	// downstream -> upstream
	go func() {
		// up/down连接断开后，2个goroutine都会err，只在一个地方清理
		defer func() {
			// 通知cleanUpBrokenConn清理
			brokenConnCh <- proxyEndMsg{
				service:  serviceName,
				endpoint: endpoint,
				id:       roomId,
			}
		}()
		defer downConn.Close()
		defer upConn.Close()
		_, err = io.Copy(upConn, downConn)
		// 不检测也可以
		if err != nil {
			return
		}
	}()
	return nil
}

func (p *proxy) cleanUpBrokenConn() {
	for cn := range brokenConnCh {
		if selector, ok := p.rt.SelectorMap[cn.service]; ok {
			selector.ReleaseEndpoint(cn.endpoint, cn.id)
		}
	}
}
