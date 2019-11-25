package internal

import (
	"fmt"
	"log"
	"testing"
)

const size = 1024

// stackCopy recursively runs increasing the size
// of the stack.
func stackCopy(s *string, c int, a [size]int) {
	println(c, s, *s)

	c++
	if c == 10 {
		return
	}

	stackCopy(s, c, a)
}

// Sample program to show how stacks grow/change. 当栈增长或者收缩时，goroutine 中的栈内存会被一块新的内存替换
func TestStack(t *testing.T) {
	s := "HELLO"
	stackCopy(&s, 0, [size]int{})
}

type user struct {
	name  string
	email string
}

//go:noinline
func createUserV1() user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V1", &u)
	// 返回值的拷贝，不同的函数栈
	return u
}

//go:noinline
func createUserV2() *user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V2", &u)
	// 返回指针，编译器逃逸分析，在堆中构造值
	return &u // 为了代码可读性，在最后&
}

func TestStack2(t *testing.T) {
	u1 := createUserV1()
	u2 := createUserV2()

	// 逃逸分析报告 go build -gcflags "-m -m"
	println("u1", &u1, "u2", &u2)
}

func TestChan(t *testing.T) {
	ch := make(chan struct{}, 1)
	close(ch)
	// ch会阻塞，在close后可以读到默认值
	s := <-ch
	s, ok := <-ch
	if ok {

	}
	log.Println(s)
}

func w(s []int )[]int {
	fmt.Printf("in w %v %#v ", s, &s[0])
	//append后s指定新复制的内存指针了，不再指向原来的内存
	s = append(s, 0)
	fmt.Printf("after append w %v %#v ", s, &s[0])
	return s
}

func TestAppend(t *testing.T){
	// https://www.cnblogs.com/sunsky303/p/11807281.html
	s:=[]int{1}
	fmt.Printf("out of w %v %#v ",s ,&s[0] )
	s2:=w(s)
	s = append(s, 1)
	fmt.Printf("after w %v %#v %v %#v ",s ,&s[0],s2 ,&s2[0] )
}