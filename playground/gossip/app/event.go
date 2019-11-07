package app

import (
	"github.com/hashicorp/memberlist"
	"log"
)

// 节点变化的事件
type event struct{}

// NotifyJoin is invoked when a node is detected to have joined.
// The Node argument must not be modified.
func (ev *event) NotifyJoin(node *memberlist.Node) {
	log.Println("A node has joined:", node.String())
}

// NotifyLeave is invoked when a node is detected to have left.
// The Node argument must not be modified.
func (ev *event) NotifyLeave(node *memberlist.Node) {
	log.Println("A node has left:", node.String())
}

// NotifyUpdate is invoked when a node is detected to have
// updated, usually involving the meta data. The Node argument
// must not be modified.
func (ev *event) NotifyUpdate(node *memberlist.Node) {
	log.Println("A node was updated:", node.String())
}
