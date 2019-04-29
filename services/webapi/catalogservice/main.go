package catalogservice

import (
	"context"
	"os"
	"os/signal"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peizhong/letsgo/framework"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

/*
dig for di
logrus for log
gin for web
gorm for mysql
redis for cache

register to service discovery
*/

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)
}

func configService() (*dig.Container, func()) {
	appsettings := framework.GetAppsettings("")

	container := dig.New()

	// appsettings
	if err := container.Provide(func() framework.Appsettings {
		return appsettings
	}); err != nil {
		log.Error(err.Error())
	}

	// config mysql
	connMySQLStr := appsettings.ConnectionStrings["MallDB"]
	dbPing, err := framework.NewMySQLConn(connMySQLStr)()
	if err != nil {
		log.WithField("MySQL connectionstring", connMySQLStr).Error(err.Error())
	}
	// test connection only
	dbPing.Close()

	// config redis
	connRedisStr := appsettings.ConnectionStrings["Redis"]
	rdPool, err := framework.InitRedisPool(connRedisStr)
	if err != nil {
		log.WithField("Redis connectionstring", connRedisStr).Error(err.Error())
	}

	// config GoContext
	container.Provide(func() *framework.GoContext {
		return &framework.GoContext{
			Context: context.Background(),
			// 只提供创建的工厂
			GetDatabase: framework.NewMySQLConn(connMySQLStr),
			GetCache:    framework.NewRedisConn(rdPool),
		}
	})

	// config controller
	container.Provide(func(c *framework.GoContext) *ProductController {
		return &ProductController{
			GoContext: c,
		}
	})

	// close connection before program exit
	return container, func() {
		log.Info("close db connection before exit")
		rdPool.Close()
	}
}

func configGin(serviceProvider *dig.Container) *gin.Engine {
	r := gin.Default()
	r.GET("/product/:id", func(gc *gin.Context) {
		err := serviceProvider.Invoke(func(c *ProductController) {
			c.HTTPContext = gc
			c.GetProduct()
		})
		if err != nil {
			log.Error(err.Error())
		}
	})
	r.GET("/products/:pageindex", func(gc *gin.Context) {
		err := serviceProvider.Invoke(func(c *ProductController) {
			c.HTTPContext = gc
			c.GetProducts()
		})
		if err != nil {
			log.Error(err.Error())
		}
	})
	ginpprof.Wrapper(r)
	return r
}

func _main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	serviceProvider, closer := configService()
	go func() {
		g := configGin(serviceProvider)
		g.Run(":8080")
	}()
	<-ch
	closer()
	log.Info("Program exit")
}

func Run() {
	g := configGin(serviceProvider)
	g.Run(":8080")
}
