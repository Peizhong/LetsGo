package main

import (
	"context"
	"fmt"
	"github.com/peizhong/letsgo/pkg/config"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/peizhong/letsgo/internal"
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
		Addr:    net.JoinHostPort("", strconv.Itoa(config.HTTPApp1Port)),
		Handler: mux,
	}
	log.Println("basic web service is on", config.HTTPApp1Port)
	if err := srv.ListenAndServeTLS(config.CertCrt, config.CertKey); err != nil {
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
