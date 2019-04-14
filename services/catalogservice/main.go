package main

import (
	"os"
	"os/signal"

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
	connStr := appsettings.ConnectionStrings["MallDB"]
	db, err := framework.NewMySQLConn(connStr)
	if err != nil {
		log.WithField("MySQL connectionstring", connStr).Error(err.Error())
	}

	// config redis
	connStr = appsettings.ConnectionStrings["Redis"]
	rdPool, err := framework.NewRedisPool(connStr)
	if err != nil {
		log.WithField("Redis connectionstring", connStr).Error(err.Error())
	}

	// config dbContext
	container.Provide(func() *framework.DbContext {
		return &framework.DbContext{
			Database: db,
			Cache:    rdPool.Get(),
		}
		// todo: 请求结束后关闭连接
	})

	// config controller
	container.Provide(func(c *framework.DbContext) *ProductController {
		return &ProductController{
			DbContext: c,
		}
	})

	// close connection before program exit
	return container, func() {
		log.Info("close db connection before exit")
		db.Close()
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
	return r
}

func main() {
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
