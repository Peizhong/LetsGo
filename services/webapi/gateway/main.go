package gateway

import (
	"net/http"
	"os"
	"os/signal"

	log "github.com/peizhong/letsgo/framework/log"

	"github.com/gorilla/mux"
)

const apiPreFix = "/api/"

/*
consul for service discovery
*/

func _main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		Run()
	}()
	select {
	case <-ch:
		log.Info("Program exit")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("I do nothing")
}

func Run() {
	r := mux.NewRouter()
	r.PathPrefix(apiPreFix).HandlerFunc(homeHandler)
	r.Use(errorMiddleware, loggingMiddleware, tracingMiddleware, reRoutingMiddleware, requestMiddleware, responseMiddleware)
	http.Handle("/", r)
	log.Info("api_gatewayservice is on")
	if err := http.ListenAndServe("localhost:8010", nil); err != nil {
		log.Error(err.Error())
	}
}
