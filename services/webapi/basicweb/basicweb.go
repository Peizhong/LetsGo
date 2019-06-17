package basicweb

import (
	"context"
	"fmt"
	"letsgo/framework/log"
	"net/http"
)

var (
	srv *http.Server
)

type BasicWeb struct {
}

type basicHandler struct {
	responseText string
}

func (th *basicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, th.responseText)
}

func (*BasicWeb) Start() {
	mux := http.NewServeMux()
	mux.Handle("/", &basicHandler{responseText: "hello world"})
	mux.Handle("/help", &basicHandler{responseText: "help me"})
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", 8001),
		Handler: mux,
	}
	log.Info("basic web service is on")
	if err := srv.ListenAndServe(); err != nil {
		log.Errorf(err.Error())
	}
}

func (*BasicWeb) Stop() {
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
