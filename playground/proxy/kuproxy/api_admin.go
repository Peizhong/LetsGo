package main

import (
	"net/http"
	"os"
)

// adminHandler: admin api handler
type adminHandler struct {
}

func (adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message := make(map[string]string)
	message["hello"] = "Hello, i will show you some api"
	message["KUBERNETES_SERVICE_HOST"] = os.Getenv("KUBERNETES_SERVICE_HOST")
	ResponseJson(w, message)
}
