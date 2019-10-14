package main

import "github.com/peizhong/letsgo/playground/gossip/cmd"

func info() string{
	msg:=`
memberlist is a Go library that manages cluster membership and member failure detection using a gossip based protocol.
all distributed systems require membership, and memberlist is a re-usable solution to managing cluster membership and node failure detection
1. gossip的原理
2. membership实现
`
	return msg
}

func main()  {
	cmd.Execute()
}
