package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuncParamMatch(t *testing.T) {
	test := func(a, b, c, d int) string {
		s := fmt.Sprint(a, b, c, d)
		return s
	}
	res := FuncParamMatch(test, 1, 2, 3, 4)
	assert.Nil(t, res)
}

func TestTemplate2(t *testing.T) {
	Template2()
}
