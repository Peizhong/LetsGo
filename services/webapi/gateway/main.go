package gateway

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/peizhong/letsgo/framework/log"

	"github.com/gorilla/mux"
)

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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", r.RequestURI)
}

func Run() {
	r := mux.NewRouter()
	r.PathPrefix("/api/").HandlerFunc(homeHandler)
	r.Use(errorMiddleware, loggingMiddleware, tracingMiddleware, reRoutingMiddleware)
	http.Handle("/", r)
	log.Info("api_gatewayservice is on")
	if err := http.ListenAndServe("localhost:8010", nil); err != nil {
		log.Error(err.Error())
	}
}
