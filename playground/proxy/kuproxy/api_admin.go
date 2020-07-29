package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
)

// adminHandler: admin api handler
type adminHandler struct {
}

func (adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message := make(map[string]string)
	message["hello"] = "Hello, i will show you some api"
	message["KUBERNETES_SERVICE_HOST"] = os.Getenv("KUBERNETES_SERVICE_HOST")
	if message["KUBERNETES_SERVICE_HOST"] != "" {
		var discovery proxy.Discovery = &proxy.K8sServiceDiscovery{}
		endpoints, _ := discovery.Endpoints("kubproxy-service")
		message["kubproxy-service"] = strings.Join(endpoints, ",")
	}
	ResponseJson(w, message)
}
