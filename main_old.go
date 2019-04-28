package main

import (
	_ "net/http/pprof"
	// load myysql driver
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/dig"

	log "github.com/sirupsen/logrus"

	"github.com/peizhong/letsgo/framework"
	"github.com/peizhong/letsgo/gonet"
	"github.com/peizhong/letsgo/gonet/middlewares"
)

// 如果导入了多个包，先初始化包的参数，然后init()，最后执行package的main()
func initOld() {
	// 将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

func configGoNetMiddleware() {
	// 请求进来时的最先遇到的在前面
	gonet.AddMiddleware(&middlewares.FirstMiddleware{})
	gonet.AddMiddleware(&middlewares.CacheMiddleware{})
	gonet.AddMiddleware(&middlewares.ReRouteMiddleware{})
	gonet.AddMiddleware(&middlewares.RequestMiddleware{})
	gonet.AddMiddleware(&middlewares.ResponseMiddleware{})
}

func mainOld() {
	//cpplib.CppTest()
	container := dig.New()
	container.Provide(func() framework.Appsettings {
		return framework.GetAppsettings("")
	})
	container.Provide(func(settings framework.Appsettings) (*middlewares.CacheMiddleware, error) {
		return &middlewares.CacheMiddleware{
			Settings: settings,
		}, nil
	})
	container.Invoke(func(md *middlewares.CacheMiddleware) {
		log.Info(md.Settings)
	})
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
	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Info("The ice breaks!", taken)
}
