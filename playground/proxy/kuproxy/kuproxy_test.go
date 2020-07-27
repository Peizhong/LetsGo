package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestK8sServiceDiscovery(t *testing.T) {
	var k8sdiscovery discovery = K8sServiceDiscovery{}
	res, err := k8sdiscovery.Endpoints("nginx-service")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	log.Println(res)
}
