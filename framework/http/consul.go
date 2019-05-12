package http

import (
	"fmt"

	"github.com/peizhong/letsgo/framework/log"
)

/*
Consul: Consul is a service mesh solution providing a full featured control plane with:
service discovery: api or mysql, DNS or HTTP
health checking
kv store: configuration
secure serivce communication: TLS
multi datacenter: segmentation functionality
//https://www.consul.io/docs/internals/architecture.html
agent: server or client mode
client: talk to server
server: elect leader
*/

func RegisterConsul(name, address string, port int, consul string) bool {
	// https://www.consul.io/api/agent/service.html#register-service
	sample := `
	{
		"Name": "redis",
		"Address": "127.0.0.1",
		"Port": 8000,
		"EnableTagOverride": false,
		"Check": {
		  "DeregisterCriticalServiceAfter": "90m",
		  "HTTP": "http://localhost:5000/health",
		  "Interval": "10s"
		},
		"Weights": {
		  "Passing": 10,
		  "Warning": 1
		}
	  }
	`
	response, err := Do("PUT", fmt.Sprintf("%v/v1/agent/service/register", consul), nil, sample)
	log.Info("%v", response.String())
	return err == nil
}
