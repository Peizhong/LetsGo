package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	res, err := Get("http://www.baidu.com/s", []Header{
		{"Content-Type", "application/json"},
		{"Connection", "keep-alive"},
	}, "wd", 132456)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}
