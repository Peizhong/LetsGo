package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTcpListen(t *testing.T) {
	// 为0时随机分配
	l, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)
	fmt.Println(l.Addr())
}

func TestRWLock(t *testing.T) {
	one := struct {
		sync.RWMutex
		value int
	}{}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		// read something
		fmt.Println("thread 1 read lock")
		one.RLock()
		fmt.Println("in thread 1, value is", one.value)
		time.Sleep(3 * time.Second)
		one.RUnlock()
		fmt.Println("thread 1 read unlock")
		wg.Done()
	}()
	time.Sleep(100)
	go func() {
		// write something
		fmt.Println("thread 2 write lock")
		one.Lock()
		fmt.Println("thread 2 get write lock")
		one.value = 1
		fmt.Println("in thread 2, set value", one.value)
		one.Unlock()
		fmt.Println("thread 2 write unlock")
		wg.Done()
	}()
	go func() {
		// read something too
		fmt.Println("thread 3 read lock")
		one.RLock()
		fmt.Println("in thread 3, value is", one.value)
		time.Sleep(3 * time.Second)
		one.RUnlock()
		fmt.Println("thread 3 read unlock")
		wg.Done()
	}()
	wg.Wait()
}
