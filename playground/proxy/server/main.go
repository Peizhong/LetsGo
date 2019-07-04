package main

import (
	"github.com/peizhong/letsgo/internal"
	"io"
	"log"
	"strconv"
	"sync"

	"fmt"
	"net"
	"time"
)

const (
	UPDATE = 1
)

/*
内网多台机器，有公网ip的服务器做穿透代理
内网机器跑客户端，公网机器跑控制器
https://zhuanlan.zhihu.com/p/70333423
https://zhuanlan.zhihu.com/p/24544885
https://github.com/fatedier/frp
*/

func start() {
	//监听控制端口8009
	go makeControl()
	//监听服务端口8007
	go makeAccept()
	//监听转发端口8008
	go makeForward()
	//定时释放连接
	go releaseConnMatch()
	//执行tcp转发
	tcpForward()
}

func stop() {

}

var (
	cache             *net.TCPConn = nil
	connListMap                    = make(map[string]*ConnMatch)
	lock                           = sync.Mutex{}
	connListMapUpdate              = make(chan int)
)

type ConnMatch struct {
	accept        *net.TCPConn //8007 tcp链路 accept
	acceptAddTime int64        //接受请求的时间
	tunnel        *net.TCPConn //8008 tcp链路 tunnel
}

// makeControl 接受代理客户端的连接
func makeControl() {
	heartBeat := func(conn *net.TCPConn) {
		for {
			//一旦有客户端连接到服务端的话，服务端每隔2秒发送hi消息给到客户端
			//如果发送不出去，则认为链路断了，清除cache连接
			_, e := conn.Write(([]byte)("hi\n"))
			if e != nil {
				cache = nil
			}
			time.Sleep(time.Second * 2)
		}
	}

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8009")
	//打开一个tcp断点监听
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("控制端口已经监听")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		fmt.Println("新的代理客户端连接到控制端服务进程:" + tcpConn.RemoteAddr().String())
		if cache != nil {
			fmt.Println("已经存在一个代理客户端连接!")
			//直接关闭掉多余的客户端请求
			tcpConn.Close()
		} else {
			cache = tcpConn
		}
		go heartBeat(tcpConn)
	}
}

// makeAccept 8007, 浏览器从外网访问这个地址，实际是访问内网的8000
func makeAccept() {
	addConnMathAccept := func(accept *net.TCPConn) {
		//加锁防止竞争读写map
		lock.Lock()
		defer lock.Unlock()
		now := time.Now().UnixNano()
		connListMap[strconv.FormatInt(now, 10)] = &ConnMatch{accept, time.Now().Unix(), nil}
	}

	sendMessage := func(message string) {
		fmt.Println("send Message " + message)
		if cache != nil {
			_, e := cache.Write([]byte(message))
			if e != nil {
				fmt.Println("消息发送异常")
				fmt.Println(e.Error())
			}
		} else {
			fmt.Println("没有代理客户端连接，无法发送消息")
		}
	}

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8007")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected 8007:" + tcpConn.RemoteAddr().String())
		addConnMathAccept(tcpConn)
		sendMessage("new\n")
	}
}

// makeForward 代理客户端连接服务器的8009后，还要再连接8008端口用于传输8007端口数据
func makeForward() {
	configConnListTunnel := func(tunnel *net.TCPConn) {
		//加锁解决竞争问题
		lock.Lock()
		used := false
		for _, connMatch := range connListMap {
			//找到tunnel为nil的而且accept不为nil的connMatch
			if connMatch.tunnel == nil && connMatch.accept != nil {
				//填充tunnel链路
				connMatch.tunnel = tunnel
				used = true
				//这里要break，是防止这条链路被赋值到多个connMatch！
				break
			}
		}
		if !used {
			//如果没有被使用的话，则说明所有的connMatch都已经配对好了，直接关闭多余的8008链路
			fmt.Println(len(connListMap))
			_ = tunnel.Close()
			fmt.Println("关闭多余的tunnel")
		}
		lock.Unlock()
		//使用channel机制来告诉另一个方法已经就绪
		connListMapUpdate <- UPDATE
	}

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8008")
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("A client connected 8008 :" + tcpConn.RemoteAddr().String())
		configConnListTunnel(tcpConn)
	}
}

func tcpForward() {
	joinConn2 := func(conn1 *net.TCPConn, conn2 *net.TCPConn) {
		f := func(local *net.TCPConn, remote *net.TCPConn) {
			//defer保证close
			defer local.Close()
			defer remote.Close()
			//使用io.Copy传输两个tcp连接，
			_, err := io.Copy(local, remote)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("join Conn2 end")
		}
		go f(conn2, conn1)
		go f(conn1, conn2)
	}

	for {
		select {
		case <-connListMapUpdate:
			lock.Lock()
			for key, connMatch := range connListMap {
				//如果两个都不为空的话，建立隧道连接
				if connMatch.tunnel != nil && connMatch.accept != nil {
					fmt.Println("建立tcpForward隧道连接")
					// 外网访问与代理客户端通讯
					go joinConn2(connMatch.accept, connMatch.tunnel)
					//从map中删除
					delete(connListMap, key)
				}
			}
			lock.Unlock()
		}
	}
}

// releaseConnMatch释放无效的代理客户端
func releaseConnMatch() {
	for {
		lock.Lock()
		for key, connMatch := range connListMap {
			//如果在指定时间内没有tunnel(代理客户端没有连接8008)的话，则释放该连接
			if connMatch.tunnel == nil && connMatch.accept != nil {
				if time.Now().Unix()-connMatch.acceptAddTime > 5 {
					fmt.Println("释放超时连接")
					err := connMatch.accept.Close()
					if err != nil {
						fmt.Println("释放连接的时候出错了:" + err.Error())
					}
					delete(connListMap, key)
				}
			}
		}
		lock.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func simpleServer() {
	hijack := "HTTP/1.1 200 OK\r\n\r\nwulala"
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		log.Println(string(buf))
		if err == nil {
			n, err = client.Write([]byte(hijack))
			client.Write([]byte{})
			log.Println(n)
		}
		err = client.Close()
		if err == nil {
			println("close client")
		}
	}
}

func main() {
	internal.Host(simpleServer, stop)
}
