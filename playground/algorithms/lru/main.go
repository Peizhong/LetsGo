package main

import (
	"container/list"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

type elem struct {
	value       interface{}
	front, next *elem
}

type milru struct {
	l  *list.List
	m  map[int]*elem
	mx sync.Mutex
}

func dolru() {
	l := &milru{
		l: list.New(),
	}
	m := l.l.PushBack("middle")
	l.l.PushBack("end")
	l.l.PushFront("front")
	// 使用了放到后面
	l.l.MoveToBack(m)
}

func main() {
	l, _ := lru.New(128)
	for i := 0; i < 256; i++ {
		l.Add(i, nil)
	}
	if l.Len() != 128 {
		panic(fmt.Sprintf("bad len: %v", l.Len()))
	}
}
