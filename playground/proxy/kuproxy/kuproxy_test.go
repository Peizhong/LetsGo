package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/peizhong/letsgo/playground/proxy/kuproxy/proxy"
	"github.com/stretchr/testify/assert"
)

func UniqueKey(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func TestK8sServiceDiscovery(t *testing.T) {
	var k8sdiscovery proxy.Discovery = &proxy.K8sServiceDiscovery{}
	res, err := k8sdiscovery.Endpoints("nginx-service")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	log.Println(res)
}

func TestSelectService(t *testing.T) {
	selector := proxy.NewSelector("Test", proxy.NewRuntime())
	roomId0 := "room0"
	endpoint0, _, err := selector.SelectEndpoint(roomId0)
	assert.NoError(t, err)
	endpoint1, _, err := selector.SelectEndpoint(roomId0)
	assert.NoError(t, err)
	assert.Equal(t, endpoint0, endpoint1)

	room1 := "room1"
	endpoint2, _, err := selector.SelectEndpoint(room1)
	assert.NoError(t, err)
	endpoint3, _, err := selector.SelectEndpoint(room1)
	assert.NoError(t, err)
	assert.Equal(t, endpoint2, endpoint3)

	assert.NotEqual(t, endpoint0, endpoint2)

	selector.ReleaseEndpoint(endpoint2, room1)

	status := selector.LoadStatus()
	assert.NotEmpty(t, status)
}

func TestRoomSerive(t *testing.T) {
	selector := proxy.NewSelector("Test", proxy.NewConfig())
	roomId0 := UniqueKey("RoomSerivce")
	endpoint0, _, _ := selector.SelectEndpoint(roomId0)
	endpoint1, _, _ := selector.SelectEndpoint(roomId0)
	status := selector.LoadStatus()
	targetRoom, ok := status[roomId0]
	assert.True(t, ok)
	assert.True(t, len(targetRoom) == 2)

	selector.ReleaseEndpoint(endpoint0, roomId0)
	selector.ReleaseEndpoint(endpoint1, roomId0)
	status = selector.LoadStatus()
	targetRoom, ok = status[roomId0]
	assert.False(t, ok)
}

func BenchmarkSelectService(b *testing.B) {
	selector := proxy.NewSelector("Test", proxy.NewConfig())
	changeRoom := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if n%2 == 1 {
			changeRoom++
		}
		roomId := fmt.Sprintf("test%v", changeRoom)
		selector.SelectEndpoint(roomId)

	}
}
