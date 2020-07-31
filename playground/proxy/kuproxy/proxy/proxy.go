package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type proxy struct {
	rt *Runtime
}

const (
	RoomIdParam      = "roomId"
	ServiceNameParam = "serviceName"
	// k8s中service的名字

	// todo：本地调试
	DefaultRoomId      = "DefaultRoomId"
	DefaultServiceName = "DefaultServiceName"
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
		panic(err)
	}
	p := &proxy{
		rt: rt,
	}
	go p.cleanUpBrokenConn()
	for {
		conn, err := l.Accept()
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
		roomId = DefaultRoomId
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
