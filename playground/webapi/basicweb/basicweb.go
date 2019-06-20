package main

import (
	"context"
	"fmt"
	"github.com/peizhong/letsgo/internal"
	"log"
	"net/http"
)

var (
	srv *http.Server
)

type basicHandler struct {
	responseText string
}

func (th *basicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, th.responseText)
}

func Start() {
	mux := http.NewServeMux()
	mux.Handle("/", &basicHandler{responseText: "hello world"})
	mux.Handle("/help", &basicHandler{responseText: "help me"})
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", 8001),
		Handler: mux,
	}
	log.Println("basic web service is on")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}

func Stop() {
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

func main() {
	internal.Host(Start, Stop)
}
