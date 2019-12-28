package main

import (
	"github.com/google/btree"
	"math"
	"testing"
	"time"
)

func TestCh(t *testing.T) {
	ch := make(chan struct{})
	v := math.Log10(20)
	println(v)
	go func() {
		select {
		case <-ch:
			println("end 1")
		}
	}()
	go func() {
		select {
		case <-ch:
			println("end 2")
		}
	}()
	ch <- struct{}{}
	<-time.After(time.Second)
}

type limitedBroadcast struct {
	transmits int   // btree-key[0]: Number of transmissions attempted.
	msgLen    int64 // btree-key[1]: copied from len(b.Message())
	id        int64 // btree-key[2]: unique incrementing id stamped at submission time

	name string // set if Broadcast is a NamedBroadcast
}

func (b *limitedBroadcast) Less(than btree.Item) bool {
	o := than.(*limitedBroadcast)
	if b.transmits < o.transmits {
		return true
	} else if b.transmits > o.transmits {
		return false
	}
	if b.msgLen > o.msgLen {
		return true
	} else if b.msgLen < o.msgLen {
		return false
	}
	return b.id > o.id
}

func TestLess(t *testing.T) {
	greaterOrEqual := &limitedBroadcast{
		transmits: 1,
		msgLen:    1000,
		id:        math.MaxInt64,
	}
	lessThan := &limitedBroadcast{
		transmits: 1,
		msgLen:    math.MaxInt64,
		id:        math.MaxInt64,
	}
	// todo: 是bug吗？greaterOrEqual没有小于lessThan
	less := greaterOrEqual.Less(lessThan)
	println(less)
}
