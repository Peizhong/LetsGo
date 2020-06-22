package main

import (
	"fmt"
	"testing"
)

// go test -test.run TestRedis
func TestRedis(t *testing.T) {
	rds := &GoRedis{}
	rds.Init()
	rds.SetString("wpz", "www")
}

// go test -bench=. -benchtime="3s" -cpuprofile profile_cpu.out
// pprof -http=:8080 profile_cpu.out

func BenchmarkDoRedis(b *testing.B) {
	var cacher Cacher
	cacher = &GoRedis{}
	cacher.Init()
	// 插入x条记录
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		_ = cacher.SetString(key, fmt.Sprint("value", i))
		_, err := cacher.GetString(key)
		if err != nil {
			// log.Println(key, "empty")
		}
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

		}
	})
	// 检查内存
}

func BenchmarkDoEtcd(b *testing.B) {
	var cacher Cacher
	cacher = &GoEtcd{}
	cacher.Init()
	// 插入x条记录
	for i := 0; i < b.N; i++ {
		key := fmt.Sprint(i)
		_ = cacher.SetString(key, fmt.Sprint("value", i))
		_, err := cacher.GetString(key)
		if err != nil {
			// log.Println(key, "empty")
		}
	}
	// 检查内存
}
