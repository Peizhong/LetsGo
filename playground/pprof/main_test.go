package main

import (
	"log"
	"testing"
	"unsafe"
)

// go test -gcflags "-N -l" -bench .

// go test -run TestPubMessage -trace=trace.out
// go tool trace trace.out

func BenchmarkNoPad_Increase(b *testing.B) {
	nopad := &NoPad{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			nopad.Increase()
		}
	})
}

func BenchmarkPad_Increase(b *testing.B) {
	pad := &Pad{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pad.Increase()
		}
	})

}

func BenchmarkMPad_Increase(b *testing.B) {
	mpad := &MPad{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mpad.Increase()
		}
	})
}

func BenchmarkSysPad_Increase(b *testing.B) {
	syspad := &SysPad{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			syspad.Increase()
		}
	})
}

func TestSize(t *testing.T) {
	nopad := NoPad{}
	pad := Pad{}
	mpad := MPad{}
	syspad := SysPad{}
	log.Println(unsafe.Sizeof(nopad))
	log.Println(unsafe.Sizeof(pad))
	log.Println(unsafe.Sizeof(mpad))
	log.Println(unsafe.Sizeof(syspad))
}
