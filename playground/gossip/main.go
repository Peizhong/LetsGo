package main

import (
	"github.com/peizhong/letsgo/playground/gossip/cmd"
	"os"
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
		os.Args = append(os.Args, []string{"run", "-n", "n0", "-p", "8081"}...)
	}
	cmd.Execute()
}
