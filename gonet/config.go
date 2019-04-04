package gonet

import (
	"encoding/json"
	"os"
	"regexp"
)

type DownStreamConfig struct {
	Host string
	Port string
}

type RouteConfig struct {
	ServiceName            string
	UpstreamPathTemplate   string
	UpstreamRegexp         *regexp.Regexp `json:",omitempty"`
	DownstreamPathTemplate string
	DownstreamHostAndPorts []DownStreamConfig
	UpstreamHttpMethod     []string
}

type GatewayConfig struct {
	Author string
	Routes []*RouteConfig
}

var gatewayConfig GatewayConfig

func LoadConfig(configFile string) {
	file, err := os.Open(configFile)
	// 函数返回时，关闭文件
	defer file.Close()
	if err != nil {
	} else {
		json.NewDecoder(file).Decode(&gatewayConfig)
	}
	// make sure route regexp is correct
	for _, route := range gatewayConfig.Routes {
		route.UpstreamRegexp = regexp.MustCompile(route.UpstreamPathTemplate)
	}
}
