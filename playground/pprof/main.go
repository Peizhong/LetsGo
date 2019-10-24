package main

import (
	"golang.org/x/sys/cpu"
	"sync/atomic"
)

type NoPad struct {
	a uint64
	b uint64
	c uint64
}

func (np *NoPad) Increase() {
	atomic.AddUint64(&np.a,1)
	atomic.AddUint64(&np.b,1)
	atomic.AddUint64(&np.c,1)
}

type Pad struct {
	a   uint64
	// 数据被多个goroutine访问，如果多个cpu同时访问某个变量，可能cpu把变量和相邻数据都读到缓存中
	// 当cpu1更新变量时，cpu2的变量无效(false sharing)，导致cpu2没用到的变量也更新
	// 将变量补足为64字节，填充一个cacheline
	_p1 [8]uint64
	b   uint64
	_p2 [8]uint64
	c   uint64
	_p3 [8]uint64
}

func (p *Pad) Increase() {
	atomic.AddUint64(&p.a,1)
	atomic.AddUint64(&p.b,1)
	atomic.AddUint64(&p.c,1)
}

type MPad struct {
	_p1 [7]uint64
	a   uint64
	_p2 [7]uint64
	b   uint64
	_p3 [7]uint64
	c   uint64
}

func (mp *MPad) Increase() {
	atomic.AddUint64(&mp.a,1)
	atomic.AddUint64(&mp.b,1)
	atomic.AddUint64(&mp.c,1)
}

type SysPad struct {
	_ cpu.CacheLinePad
	a uint64
	_ cpu.CacheLinePad
	b uint64
	_ cpu.CacheLinePad
	c uint64
}

func (sp *SysPad) Increase() {
	atomic.AddUint64(&sp.a,1)
	atomic.AddUint64(&sp.b,1)
	atomic.AddUint64(&sp.c,1)
}