package main

// discovery upstream
type discovery interface {
	// endpoints 返回endpoints信息
	Endpoints() []string
}
