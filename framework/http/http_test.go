package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	res, err := Do("GET", "http://www.baidu.com/s", []Header{
		{"Content-Type", "application/json"},
		{"Connection", "keep-alive"},
	}, "", "wd", 132456)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func Test_Consul(t *testing.T) {
	res := RegisterConsul("", "", 8000, "http://127.0.0.1:8500")
	assert.True(t, res)
}
