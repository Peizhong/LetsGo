package proxy

type Runtime struct {
	Discovery   Discovery
	SelectorMap map[string]Selector
	Config      Config
}

func NewRuntime() *Runtime {
	rt := &Runtime{}
	rt.Discovery = NewSerivceDiscovery()
	rt.SelectorMap = make(map[string]Selector)
	rt.Config = NewConfig()
	return rt
}
