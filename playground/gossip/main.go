package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/peizhong/letsgo/playground/gossip/cmd"
)

var (
	// go build -ldflags "-X 'main.Version=0.0.1'"
	Version string
	// pprof
	// go tool pprof -http=:8080 -no_browser gossip.exe cpu.prof
	cpuprofile = "cpu.prof"
	memprofile = "memprofile"
)

func info() string {
	msg := `
memberlist is a Go library that manages cluster membership and member failure detection using a gossip based protocol.
all distributed systems require membership, and memberlist is a re-usable solution to managing cluster membership and node failure detection
1. gossip的原理
2. membership实现
`
	return msg
}

func main() {
	if cpuprofile != "" {
		f, err := os.OpenFile(cpuprofile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
		mf, err := os.OpenFile(memprofile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("could not create Mem profile: ", err)
		}
		defer mf.Close()
		if err := pprof.WriteHeapProfile(mf); err != nil {
			log.Fatal("could not start Mem profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if len(os.Args) < 2 {
		// 默认启动
		os.Args = append(os.Args, []string{"run", "-n", "node0", "-p", "8081"}...)
	}
	println("version:", Version)
	cmd.Execute()
	println("bye")
}
