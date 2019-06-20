package main

import (
	"github.com/micro/mdns"
	"github.com/peizhong/letsgo/internal"
	"os"
)

var server *mdns.Server

func Start() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, _ := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ = mdns.NewServer(&mdns.Config{Zone: service})
}

func Stop() {
	server.Shutdown()
}

func main() {
	internal.Host(Start, Stop)
}
