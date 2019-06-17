package gateway

import (
	"context"
	"fmt"
	"letsgo/framework/config"
	"letsgo/framework/log"
	"net/http"

	"github.com/gorilla/mux"
)

const apiPreFix = "/api/"

var (
	srv *http.Server
)

/*
consul for service discovery
*/

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("I do nothing")
}

type GatewayService struct {
}

func (*GatewayService) Start() {
	r := mux.NewRouter()
	r.PathPrefix(apiPreFix).HandlerFunc(homeHandler)
	r.Use(errorMiddleware, loggingMiddleware, tracingMiddleware, reRoutingMiddleware, requestMiddleware, responseMiddleware)
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GatewayPort),
		Handler: r,
	}
	log.Info("api gateway service is on")
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf(err.Error())
	}
}

func (*GatewayService) Stop() {
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
