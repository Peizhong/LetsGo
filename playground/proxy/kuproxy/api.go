package main

import (
	"encoding/json"
	"net/http"
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
