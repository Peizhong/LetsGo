package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	// https://mikespook.com/2013/07/%E7%BF%BB%E8%AF%91go-%E7%9A%84%E8%B0%83%E5%BA%A6%E5%99%A8/
	runtime.GOMAXPROCS(2)
	var i = 0
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for n := 0; n < 10; n++ {
			i++
			fmt.Println("p1", i)
		}
		wg.Done()
	}()
	go func() {
		for n := 0; n < 10; n++ {
			i++
			fmt.Println("p2", i)
		}
		wg.Done()
	}()
	wg.Wait()
}

func example(slice []string, str string, i int) {
	panic("Want stack trace")
}
