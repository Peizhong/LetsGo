package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/peizhong/letsgo/gonet"
	"github.com/peizhong/letsgo/gonet/middlewares"
)

// 如果导入了多个包，先初始化包的参数，然后init()，最后执行package的main()
func init() {
	// 将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

func configGoNetMiddleware() {
	// 请求进来时的最先遇到的在前面
	gonet.AddMiddleware(&middlewares.FirstMiddleware{})
	gonet.AddMiddleware(&middlewares.ReRouteMiddleware{})
	gonet.AddMiddleware(&middlewares.RequestMiddleware{})
	gonet.AddMiddleware(&middlewares.ResponseMiddleware{})
}

func main() {
	//cpplib.CppTest()

	begin := time.Now()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		configGoNetMiddleware()
		gonet.LoadConfig("config/gateway.json")
		gonet.RunHTTPServer("localhost", 8080)
	}()
	<-c
	taken := time.Since(begin)
	fmt.Println("program exit, last for", taken)
}
