package gonet

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
}

var once sync.Once

var entryPoint func(*Context) error

func gatewayHandler(w http.ResponseWriter, req *http.Request) {
	once.Do(func() {
		entryPoint = BuildPipeline()
	})
	context := &Context{
		SrcPath:   req.RequestURI,
		Responser: w,
	}
	err := entryPoint(context)
	if err != nil {
		fmt.Println("request error")
	}
}

// RunHTTPServer start gonet http server
func RunHTTPServer(ip string, port int) (err error) {
	address := fmt.Sprintf("%v:%v", ip, port)
	http.HandleFunc("/", gatewayHandler)
	fmt.Println("HTTP server listening at:", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
	return
}
