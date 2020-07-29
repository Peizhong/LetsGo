package proxy

import (
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

const (
	RoomIdParam      = "roomId"
	ServiceNameParam = "serviceName"
	// k8s中service的名字
	// todo: 从环境变量读取
	DefaultServiceName = "xsdj-vrol"
)

type proxyEndMsg struct {
	service, endpoint, id string
}

var (
	ParseRequestError  = errors.New("Parse request error")
	serviceSelectorMap = make(map[string]selector)
	// 客户端连接断开后，通知清理
	brokenConnCh = make(chan proxyEndMsg, 10)
)

func parseRequestLine(line string) (method, requestURI, roomId, serviceName string) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	method = line[:s1]
	s2 += s1 + 1
	requestURI = line[s1+1 : s2]
	parseParams := func(param string) (target string) {
		if targetStart := strings.Index(requestURI, param); targetStart >= 0 {
			targetStringGreedy := requestURI[targetStart:]
			if targetEnd := strings.Index(targetStringGreedy, "&"); targetEnd > 0 {
				target = targetStringGreedy[:targetEnd]
			}
		}
		return
	}
	roomId = parseParams(RoomIdParam)
	// todo: 从环境变量读取
	serviceName = parseParams(ServiceNameParam)
	return
}

// todo: 换成select/epoll?
func makeTunnel(down net.Conn) error {
	buf := make([]byte, 1024)
	n, err := down.Read(buf)
	if err != nil {
		return err
	}
	// 检查请求是否是包含roomId的http请求
	_, url, roomId, serviceName := parseRequestLine(string(buf))
	if url == "" || roomId == "" {
		return ParseRequestError
	}
	if serviceName == "" {
		serviceName = DefaultServiceName
	}
	log.Println(url, roomId)
	selector, ok := serviceSelectorMap[serviceName]
	if !ok {
		selector = NewSelector(serviceName)
		serviceSelectorMap[serviceName] = selector
	}
	endpoint, err := selector.SelectEndpoint(roomId)
	if err != nil {
		return err
	}
	upconn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return err
	}
	// 将截取的tcp写回upstream
	_, err = upconn.Write(buf[:n])
	if err != nil {
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

func cleanUpBrokenConn() {
	for cn := range brokenConnCh {
		if selector, ok := serviceSelectorMap[cn.service]; ok {
			selector.ReleaseEndpoint(cn.endpoint, cn.id)
			// todo: 清理全部数量为0的房间记录
		}
	}
}
