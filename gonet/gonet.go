package gonet

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
)

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		fmt.Printf("Received data: %v", string(buf[:len]))
		conn.Write([]byte("ok"))
	}
}

// RunTCPServer start gonet server
func RunTCPServer(ip string, port int) error {
	address := fmt.Sprintf("%v:%v", ip, port)
	fmt.Println("Starting the server on: ", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return err
	}
	// 监听并接受来自客户端的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return err
		}
		go doServerStuff(conn)
	}
}

func helloServer(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("Inside HelloServer handler")
	fmt.Fprintf(w, "Hello,"+req.URL.Path[1:])
}

// RunHTTPServer start gonet http server
func RunHTTPServer(ip string, port int) (err error) {
	address := fmt.Sprintf("%v:%v", ip, port)
	http.HandleFunc("/", helloServer)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	return
}
