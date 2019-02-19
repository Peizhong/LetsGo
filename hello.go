package main

import (
	"fmt"
)

func main() {
	// 定义变量
	var b bool
	// 定义并赋值
	var n = 1
	// 简短声明
	_, v1, v2, v3 := 0, 1, 3, 5
	// 常量
	const value0 = 1234
	value1 := 1
	_, value2 := 2, value1
	// 内置数据类型
	var enable = true
	// 浮点
	num := 1.1
	var helloString = "how you doing"
	// 字符串转换为[]byte
	c := []byte(helloString)
	s2 := string(c)
	// 字符串切片
	s3 := "h" + s2[1:]
	// 没有支付转义
	sraw := s3 + `ads
	ads
	`
	// array
	nums := [4]int{1, 3, 5, 7}
	// 省略长度
	cnums := [...]int{2, 4, 6, 8}
	// 二维数组 row*column
	doubleArray := [...][4]int{{1, 2, 3, 4}, {5, 6, 7, 8}}
	// slice 不固定长度
	slice := []byte{'a', 'b', 'c', 'd'}
	// iota枚举， const中重置，每行加1
	const (
		red   = iota //0
		green = iota //1
		blue  = iota //2
	)
	fmt.Println("%s", sraw)
	fmt.Println(b, n, v1, v2, v3, value2, enable, num, nums, cnums, doubleArray, slice)
	fmt.Println("!oG ,olleH")
	// 大写变量/函数是可导出的，其他包可以读取，小写的是私有
}
