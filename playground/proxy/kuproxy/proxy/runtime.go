package proxy

import (
	"log"
	"os"
)

type Runtime struct {
	Discovery   Discovery
	SelectorMap map[string]Selector
	Config      Config
}

func NewRuntime() *Runtime {
	rt := &Runtime{}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		log.Println("in k8s")
		rt.Discovery = &K8sServiceDiscovery{}
	} else {
		rt.Discovery = &MockServiceDiscovery{}
		log.Println("in dev")
	}
	rt.SelectorMap = make(map[string]Selector)
	rt.Config = NewConfig()
	return rt
}
