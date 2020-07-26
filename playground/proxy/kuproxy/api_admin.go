package main

import (
	"net/http"
)

// adminHandler: admin api handler
type adminHandler struct {
}

func (adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, "Hello, i will show you some api")
}
