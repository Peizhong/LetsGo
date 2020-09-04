package main

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/stretchr/testify/assert"
)

func TestGlob(t *testing.T) {
	g := glob.MustCompile("a:*", ':')
	r := g.Match("a:b")
	assert.True(t, r)
	r = g.Match("a:a:c")
	assert.False(t, r)
	// superwild card **, 贪婪匹配, 不受分隔符约束
	sg := glob.MustCompile("a:**", ':')
	r = sg.Match("a:b")
	assert.True(t, r)
	r = sg.Match("a:a:c")
	assert.True(t, r)
	// domain 范围scpoe(团队->二级部门->一级部门->公司->系统):id
	//			找到所有从上级继承的domain
	// resource 应用:版本
	// action 模块:功能:操作
}
