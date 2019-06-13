package learn

import (
	"github.com/stretchr/testify/assert"
	"letsgo/framework"
	"letsgo/framework/log"
	"testing"
)

func TestBucket(t *testing.T) {
	v := make(map[int64]int8, 10)
	v[1] = 1
	v[2] = 2
	log.Info("v is %v", len(v))
	framework.WhatIsThis(v, v[1], v[2])
	if a, ok := v[123]; ok {
		assert.True(t, a > 0)
	}
}
