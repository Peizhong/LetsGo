package main

import (
	"log"
	"testing"

	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
	"github.com/stretchr/testify/assert"
)

func TestK8sServiceDiscovery(t *testing.T) {
	var k8sdiscovery proxy.Discovery = &proxy.K8sServiceDiscovery{}
	res, err := k8sdiscovery.Endpoints("nginx-service")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	log.Println(res)
}

func TestSelectService(t *testing.T) {
	selector := proxy.NewSelector("Test")
	roomId0 := "room0"
	endpoint0, err := selector.SelectEndpoint(roomId0)
	assert.NoError(t, err)
	endpoint1, err := selector.SelectEndpoint(roomId0)
	assert.NoError(t, err)
	assert.Equal(t, endpoint0, endpoint1)

	room1 := "room1"
	endpoint2, err := selector.SelectEndpoint(room1)
	assert.NoError(t, err)
	endpoint3, err := selector.SelectEndpoint(room1)
	assert.NoError(t, err)
	assert.Equal(t, endpoint2, endpoint3)

	assert.NotEqual(t, endpoint0, endpoint2)
}
