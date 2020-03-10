package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/peizhong/letsgo/internal"
	"log"
	"net/http"
	"strings"
)

type httpRedirect struct {
	server *http.Server
}

func (h *httpRedirect) start() {
	prepareGeo()

	r := mux.NewRouter()
	r.HandleFunc("/visitor", VisitorHandler)
	r.PathPrefix("/").HandlerFunc(RedirectHandler)
	h.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", IP, GatewayPort),
		Handler: r,
	}
	// 1024以下端口，user组不能访问
	log.Println("start on", GatewayPort)
	if err := h.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
func (h *httpRedirect) stop() {
	log.Println("close server")
	err := h.server.Close()
	internal.CheckError(err, "Close Server")

	log.Println("exit geo")
	err = unloadGeo()
	internal.CheckError(err, "Unload Geo")
}

// 除了特定的路径，其他全部转移到8000端口
func fixUrl(host, url string) string {
	i := strings.LastIndex(host, ":")
	if i > 0 {
		host = host[:i]
	}
	r := fmt.Sprintf("http://%s:%d%s", host, RealPort, url)
	log.Println(host, url, r)
	return r
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := fixUrl(r.Host, r.RequestURI)
	recordGeo(r.RemoteAddr)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
