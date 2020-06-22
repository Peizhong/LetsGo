package main

import (
	"log"
	"testing"
)

func dfs64(current []string, str []string, ref []bool) {
	for i := range ref {
		if ref[i] == false {
			// 且与前面的值不重复
			for ck := 0; ck < i; ck++ {
				if ref[ck] == true && str[ck] == str[i] {
					goto next
				}
			}
			bak := current[:]
			current = append(current, str[i])
			ref[i] = true
			dfs64(current, str, ref)
			current = bak
			ref[i] = false
		}
	next:
	}
	if len(current) == len(str) {
		log.Println(current)
	}
}

// go test -test.run Test46
func Test46(t *testing.T) {
	// 输入字符串，字符有重复，打印所有不重复的排列
	str := []string{"d", "d", "a", "a"}
	used := make([]bool, len(str))
	dfs64([]string{}, str, used)
}

func TestTextMatch(t *testing.T) {
	// match可以为1,0,?, 问号匹配1或0
	// 查找str中，符合条件的不重复字符串数量
	str := "00010001"
	match := "??"
	// 从0开始比较
	slen := len(match)
	res := make(map[string]struct{})
	for i := 0; i < len(str)-slen; i++ {
		ok := true
		for j := 0; j < slen; j++ {
			if match[j] != '?' && match[j] != str[i+j] {
				ok = false
				break
			}
		}
		if ok {
			res[str[i:i+slen]] = struct{}{}
		}
	}
	log.Println(len(res))
}

func TestTextNearBy(t *testing.T) {
	// 输入字符串，输出字符串中相邻字符所有组合
	// 输出去重，相同长度的按字典序排序
	str := "abc"
	slen := len(str)
	for step := 1; step < slen; step++ {
		// 取出来，再排序。。
	}
}
