package main

import (
	"github.com/peizhong/letsgo/pkg/config"
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8080", config.CertCrt, config.CertKey, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
