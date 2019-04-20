package framework

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

func IntTryParse(s string) (n int, b bool) {
	if num, err := strconv.Atoi(s); err == nil {
		return num, true
	}
	return
}

func Int64TryParse(s string) (n int64, b bool) {
	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num, true
	}
	return
}

func WhatIsThis(values ...interface{}) {
	for _, v := range values {
		spew.Dump(v)
	}
}
