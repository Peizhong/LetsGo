package main

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sync/singleflight"
)

func usecond() {
	cond := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 1; i++ {
		tt := i
		go func() {
			cond.L.Lock()
			cond.Wait()
			println(tt)
			cond.L.Unlock()
		}()
	}
	println("wait 1 second to broadcast")
	<-time.After(time.Second)
	// not required cond.L.Lock()
	// cond.L.Lock()
	cond.Signal()
	cond.Broadcast()
}

func usesingleflight() {
	flight := &singleflight.Group{}
	var wg sync.WaitGroup
	var acutally int
	keys := []string{"a", "b", "c", "a", "b", "c"}
	for i := range keys {
		wg.Add(1)
		k := keys[i]
		go func() {
			// 上一个同样的key还没结束时
			flight.Do(k, func() (interface{}, error) {
				<-time.After(time.Second)
				acutally++
				return nil, nil
			})
			wg.Done()
		}()
	}
	wg.Wait()
	println("actually", acutally)
}

func usePool() {
	type ss struct {
		v byte
	}
	pool := sync.Pool{New: func() interface{} {
		return ss{}
	}}
	v := pool.New()
	pool.Put(v)
	_ = pool.Get()
}

func syncMap() {
	var mp sync.Map
	mp.Store("aa", "bb")
	if v, ok := mp.Load("aa"); ok {
		println(v)
	}
	mp.Store(12, "34")
	mp.Delete("aa")
	mp.Range(func(k, v interface{}) bool {
		return true
	})
}

func UseAtomic() {
	type v struct {
		a, b int
	}
	av := atomic.Value{}
	av.Store(v{1, 2})
	xv := av.Load().(v)
	p := (*v)(unsafe.Pointer(uintptr(unsafe.Pointer(&xv)) + 10))
	println(xv.a, xv.b, p)
}

func UseSync() {
	usesingleflight()
	usecond()
}
