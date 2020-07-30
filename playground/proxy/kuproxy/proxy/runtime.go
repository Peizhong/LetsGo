package proxy

import "os"

type Runtime struct {
	Discovery   Discovery
	SelectorMap map[string]Selector
	Config      Config
}

func NewRuntime() *Runtime {
	rt := &Runtime{}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		rt.Discovery = &K8sServiceDiscovery{}
	} else {
		rt.Discovery = &MockServiceDiscovery{}
	}
	rt.SelectorMap = make(map[string]Selector)
	rt.Config = NewConfig()
	return rt
}
