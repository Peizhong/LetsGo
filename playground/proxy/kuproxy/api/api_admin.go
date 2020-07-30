package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
)

// AdminHandler: admin api handler
type AdminHandler struct {
	rt *proxy.Runtime
}

func NewAdminHandler(rt *proxy.Runtime) *AdminHandler {
	return &AdminHandler{
		rt: rt,
	}
}

// adminHandler: admin api handler
func (h *AdminHandler) root(w http.ResponseWriter, r *http.Request) {
	message := make(map[string]string)
	message["hello"] = "Hello, i will show you some api"
	ResponseJson(w, message)
}

func (h *AdminHandler) serviceEndpoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if serviceName := vars["serviceName"]; serviceName != "" {
		endpoints, err := h.rt.Discovery.Endpoints(serviceName)
		if err != nil {
			ResponseError(w, err.Error(), 500)
			return
		}
		ResponseJson(w, fmt.Sprintf("%s:%v", serviceName, endpoints))
		return
	}
	ResponseError(w, "no service name", 404)
}
