package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peizhong/letsgo/internal"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/visitor", VisitorHandler)
	r.PathPrefix("/").HandlerFunc(RedirectHandler)
	srv := http.Server{
		Addr:    "192.168.3.34:8080",
		Handler: r,
	}

	internal.Host(func() {
		prepareGeo()
		// 1024一下端口，user组不能访问
		log.Println("start on", ":8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}, func() {
		log.Println("close server")
		err := srv.Close()
		internal.CheckError(err, "Close Server")
		log.Println("exit geo")
		err = unloadGeo()
		internal.CheckError(err, "Unload Geo")
	})
}
