package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"testing"
	"time"

	// 宿主的cpu和容器实际cpu
	_ "go.uber.org/automaxprocs"
)

// go build -gcflags -m

//go:noinline
func simple(n int) {
	// inline: 将函数调用替换为函数主体
	// pros: 函数调用有固定的开销；栈和抢占检查。减少函数调用开销。更大的函数可能编译优化。只对叶子函数有效
	// cons: 重复代码，降低(cpu)缓存命中
	fmt.Println(n)
}

func other() {
	// nosplit: 跳过栈溢出检测
	// noescape: 不做逃逸分析。指示一个没有主体的函数
}

func oneP() {
	runtime.GOMAXPROCS(1)

	go func() {
		simple(0)
	}()

	go simple(1)

	go simple(2)

	select {}
}

var sum int

//go:norace
func add() {
	// 跳过竞态检测
	// go run -race, go build -race
	sum++
}

// https://mikespook.com/2013/07/%E7%BF%BB%E8%AF%91go-%E7%9A%84%E8%B0%83%E5%BA%A6%E5%99%A8/
func race() {
	go add()
	go add()
}

func iocopy() {
	// https://blog.go-zh.org/race-detector
	// If the given Writer implements a ReadFrom method, call writer.ReadFrom(reader)
	// Discard has an internal buffer that is shared
	// race detector flagged this racy. a non-racy version when the race detector is running.
	io.Copy(ioutil.Discard, bytes.NewReader([]byte{}))
}

func doDhannel() {

}

func cal(index string, a, b int) int {
	ret := a + b
	println(index, a, b, ret)
	return ret
}

func TestCal1(t *testing.T) {
	a := 1
	b := 2
	defer cal("1", a, cal("10", a, b))
	a = 0
	defer cal("2", a, cal("20", a, b))
	b = 1
}
func TestCal2(t *testing.T) {
	a := 1
	b := 2
	defer func() {
		cal("1", a, cal("10", a, b))
	}()
	a = 0
	defer func() {
		cal("2", a, cal("20", a, b))
	}()
	b = 1
}

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func main() {
	if false {

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				var r int
				for v := 0; v < 100000; v++ {
					for w := 0; w < 10000; w++ {
						r = v*w - v + w
					}
				}
				if r > 0 {

				}
				wg.Done()
			}()
		}
		// runtime.Gosched()
		for i := 0; i < 6; i++ {
			// 这里只有6个，会出现idleprocs，但没有自旋?
			wg.Add(1)
			go func() {
				var r int
				for v := 0; v < 100000; v++ {
					for w := 0; w < 10000; w++ {
						r = v*w - v + w
						r = v*w - v + w
					}
				}
				if r > 0 {

				}
				wg.Done()
			}()
		}
		go func() {
			<-time.After(time.Minute)
			wg.Done()
		}()
		wg.Wait()
	}
}
