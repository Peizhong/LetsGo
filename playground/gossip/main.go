package main

import (
	"github.com/peizhong/letsgo/playground/gossip/cmd"
	"os"
)

var (
	// go build -ldflags "-X 'main.Version=0.0.1'"
	Version string
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
	if len(os.Args) < 2 {
		// 默认启动
		os.Args = append(os.Args, []string{"run", "-n", "node1", "-p", "8081", "-j", "localhost:8080"}...)
	}
	println("version:", Version)
	cmd.Execute()
}
