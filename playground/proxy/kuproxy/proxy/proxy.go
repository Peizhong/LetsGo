package proxy

import (
	"errors"
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
		return err
	}
	p := &proxy{
		rt: rt,
	}
	go p.cleanUpBrokenConn()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		if err = p.makeTunnel(conn); err != nil {
			log.Println("abort tcp conn", err.Error())
			conn.Close()
		}
	}
}

// makeTunnel: 建立双向转发通道
// todo: 每个连接会产生2个goroutine，换成select/epoll?
func (p *proxy) makeTunnel(down net.Conn) error {
	buf := make([]byte, 1024)
	n, err := down.Read(buf)
	if err != nil {
		return err
	}
	// 检查请求是否是包含roomId的http请求
	_, url, roomId, serviceName := parseRequestLine(string(buf))
	if url == "" {
		log.Println("unknow request")
		return ParseRequestError
	}
	if roomId == "" {
		roomId = DefaultRoomId
	}
	if serviceName == "" {
		// 调试用，应该在http请求中包含
		serviceName = DefaultServiceName
	}
	log.Println(url, roomId, serviceName)
	selector, ok := p.rt.SelectorMap[serviceName]
	if !ok {
		selector = NewSelector(serviceName, p.rt.Config)
		p.rt.SelectorMap[serviceName] = selector
	}
	endpoint, err := selector.SelectEndpoint(roomId)
	if err != nil {
		return err
	}
	log.Println("dial to ", endpoint)
	upconn, err := net.Dial("tcp", endpoint)
	if err != nil {
		// dial失败，释放刚申请的房间
		brokenConnCh <- proxyEndMsg{
			service:  serviceName,
			endpoint: endpoint,
			id:       roomId,
		}
		return err
	}
	// 将截取的tcp写回upstream
	_, err = upconn.Write(buf[:n])
	if err != nil {
		// 写入失败，释放刚申请的房间
		brokenConnCh <- proxyEndMsg{
			service:  serviceName,
			endpoint: endpoint,
			id:       roomId,
		}
		return err
	}
	// upstream -> downstream
	go func() {
		defer down.Close()
		defer upconn.Close()
		_, err := io.Copy(down, upconn)
		if err != nil {
			return
		}
	}()
	// downstream -> upstream
	go func() {
		// 连接断开后，通知cleanUpBrokenConn清理
		// 2个goroutine，有一个err，另一个也会err，只触发一次
		defer func() {
			brokenConnCh <- proxyEndMsg{
				service:  serviceName,
				endpoint: endpoint,
				id:       roomId,
			}
		}()
		defer down.Close()
		defer upconn.Close()
		_, err = io.Copy(upconn, down)
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
