package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/peizhong/letsgo/playground/webapi/elastic/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.elastic.co/apm/module/apmgin"
	"go.elastic.co/apm/module/apmgoredis"
	"go.elastic.co/apm/module/apmhttp"

	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/mysql"
)

// https://www.elastic.co/guide/en/apm/agent/go/current/builtin-modules.html

func CheckError(err error, fmtStr string, args ...interface{}) (ce error) {
	defer func() {
		if ce != nil {
			log.Println(ce.Error())
		}
	}()
	if err == nil {
		return
	}
	if fmtStr == "" {
		ce = fmt.Errorf("err: %s", err.Error())
	}
	if len(args) == 0 {
		ce = fmt.Errorf("%s err: %s", fmtStr, err.Error())
	}
	ce = fmt.Errorf("%s err: %s", fmt.Sprintf(fmtStr, args...), err.Error())
	return
}

func connStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", config.DB().UserName, config.DB().Password, config.DB().Address, config.DB().Database)
}

type PingHandler struct {
}

func (PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

type PostDemoData struct {
	Id, Value string
}

func buildHandler() http.Handler {
	engine := gin.New()
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	var myHandler http.Handler = PingHandler{}
	tracedHandler := apmhttp.Wrap(myHandler)
	engine.GET("/ping", gin.WrapH(tracedHandler))
	adminGroup := engine.Group("/admin", func(c *gin.Context) {
		if auth := c.Request.Header.Get("Authorization"); auth != "foobar" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(adminGroup, "pprof")

	g := engine.Group("/api")
	g.Use(apmgin.Middleware(engine))
	g.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	g.GET("/host", func(c *gin.Context) {
		c.JSON(http.StatusOK, os.Hostname)
	})
	g.GET("/app/:appid", func(c *gin.Context) {
		all := struct {
			AppId string `uri:"appid"`
			Key   string `form:"key"`
			Value string
		}{}
		c.BindUri(&all)
		c.BindQuery(&all)
		c.JSON(http.StatusOK, all)
	})
	g.POST("/app/:appid", func(c *gin.Context) {
		all := struct {
			AppId string `uri:"appid"`
			Key   string `form:"key"`
			Id    string
			Value string
		}{}
		c.BindUri(&all)
		c.BindQuery(&all)
		c.BindJSON(&all)
		c.JSON(http.StatusOK, all)
	})
	return engine
}

func main() {
	dbConnectionString := connStr()
	db, err := apmgorm.Open("mysql", dbConnectionString)
	if ce := CheckError(err, ""); ce != nil {
		return
	}
	err = db.DB().Ping()
	if ce := CheckError(err, "mysql.ping"); ce != nil {
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "111111",
		DB:       0,
	})
	rclient := apmgoredis.Wrap(redisClient).WithContext(context.Background())
	err = rclient.Ping().Err()
	if ce := CheckError(err, "redis.ping"); ce != nil {
		return
	}
	engine := buildHandler()
	log.Println("listen at", config.HTTP().LitenAddress)
	if err = http.ListenAndServe(config.HTTP().LitenAddress, engine); err != nil {
		panic(err)
	}
}
