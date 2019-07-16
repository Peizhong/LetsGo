package main

import (
	"fmt"
)

func main() {
	// 创建一个容量和长度均为6的slice
	slice1 := []int{5, 23, 10, 2, 61, 33}
	// 对slices1进行切片，长度为2容量为4
	slice2 := slice1[1:3]
	fmt.Println("cap", cap(slice2))
	fmt.Println("slice2", slice2)
}
