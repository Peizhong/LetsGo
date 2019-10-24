package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
)

// go build -gcflags -m

//go:noinline
func simple(n int) {
	// inline: 将函数调用替换为函数主体
	// pros: 函数调用有固定的开销；栈和抢占检查。减少函数调用开销。更大的函数可能编译优化。只对叶子函数有效
	// cons: 重复代码，降低(cpu)缓存命中
	fmt.Println(n)
}

func other(){
	// nosplit: 跳过栈溢出检测
	// noescape: 不做逃逸分析。指示一个没有主体的函数
}

func oneP() {
	runtime.GOMAXPROCS(1)

	go func() {
		simple(0)
		for {

		}
	}()

	go simple(1)

	go simple(2)

	select {}
}

var sum int

//go:norace
func add(){
	// 跳过竞态检测
	// go run -race, go build -race
	sum++
}

func race(){
	go add()
	go add()
}

func copy(){
	io.Copy(ioutil.Discard,)
}

func main() {
	// https://mikespook.com/2013/07/%E7%BF%BB%E8%AF%91go-%E7%9A%84%E8%B0%83%E5%BA%A6%E5%99%A8/
	race()
}