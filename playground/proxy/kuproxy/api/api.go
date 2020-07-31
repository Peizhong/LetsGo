package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ResponseJson(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	bytes, err := json.Marshal(obj)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(bytes)
}

func ResponseError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func StartHttp(httpAddr string, rt *proxy.Runtime) {
	r := mux.NewRouter()
	// todo: 如果部署多个副本就用不了了
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.Handle("/metrics", promhttp.Handler())

	admin := r.PathPrefix("/admin").Subrouter()
	adminHandler := NewAdminHandler(rt)
	admin.HandleFunc("/", adminHandler.root)
	admin.HandleFunc("/services/{serviceName}", adminHandler.serviceEndpoints)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.Use(LoggingMiddleware)

	srv := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}
	log.Println("http listen at", httpAddr)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
