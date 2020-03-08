package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 除了特定的路径，其他全部转移到8000端口
const (
	hostPort  = ""
	redirPort = "8000"
)

func fixHost(host, url string) string {
	i := strings.LastIndex(host, ":")
	if i > 0 {
		host = host[:i]
	}
	r := fmt.Sprintf("http://%s:%s%s", host, redirPort, url)
	log.Println(host, url, r)
	return r
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	fixUrl := fixHost(r.Host, r.RequestURI)
	recordGeo(r.RemoteAddr)
	http.Redirect(w, r, fixUrl, http.StatusSeeOther)
}
