package app

import (
	"context"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"log"
	"time"
)

func NewRaftStorage(id uint64, peerIds []uint) {
	storage := raft.NewMemoryStorage()
	c := &raft.Config{
		ID:              id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	var peers []raft.Peer
	for _, p := range peerIds {
		peers = append(peers, raft.Peer{ID: uint64(p)})
	}
	<-time.After(time.Duration(id) * time.Second)
	log.Println(peerIds)
	node := raft.StartNode(c, peers)
	if err := node.ProposeConfChange(context.Background(), raftpb.ConfChange{
		ID:      0,
		Type:    0,
		NodeID:  0,
		Context: nil,
	}); err != nil {
		log.Println(err)
	}
}
