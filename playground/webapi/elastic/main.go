package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.elastic.co/apm/module/apmgin"
	"go.elastic.co/apm/module/apmgoredis"
	"go.elastic.co/apm/module/apmhttp"

	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/mysql"
)

// https://www.elastic.co/guide/en/apm/agent/go/current/builtin-modules.html

func main() {
	db, err := apmgorm.Open("mysql", "")
	if err == nil {
		db.DB().Ping()
	}
	var redisClient *redis.Client
	rclient := apmgoredis.Wrap(redisClient).WithContext(context.Background())
	rclient.Ping()

	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))
	var myHandler http.Handler
	tracedHandler := apmhttp.Wrap(myHandler)
	http.ListenAndServe(":8080", engine)
	http.ListenAndServe(":8081", tracedHandler)

}
