package gonet

import (
	"fmt"
	"net/http"
)

type Gateway struct {
}

// ServeHTTP handle all request, disabled other plugins
func (Gateway) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	once.Do(func() {
		entryPoint = BuildPipeline()
	})
	context := &Context{
		SrcPath:   req.RequestURI,
		Responser: w,
	}
	err := entryPoint(context)
	if err != nil {
		fmt.Println("request error")
	}
}
