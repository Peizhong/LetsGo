package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/coreos/etcd/etcdserver/stats"
	"github.com/coreos/etcd/pkg/types"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/rafthttp"
)

type raftNode struct {
	id        int // client ID for raft session
	port      int
	peers     []string // raft peer URLs
	node      raft.Node
	transport *rafthttp.Transport
}

func (rc *raftNode) Process(ctx context.Context, m raftpb.Message) error {
	return rc.node.Step(ctx, m)
}
func (rc *raftNode) IsIDRemoved(id uint64) bool                           { return false }
func (rc *raftNode) ReportUnreachable(id uint64)                          {}
func (rc *raftNode) ReportSnapshot(id uint64, status raft.SnapshotStatus) {}

func (rc *raftNode) serveRaft() {
	url, err := url.Parse(fmt.Sprintf(":%d", rc.port))
	if err != nil {
		log.Fatalf("raftexample: Failed parsing URL (%v)", err)
	}

	ln, err := net.Listen("tcp", url.Host)
	if err != nil {
		log.Fatalf("raftexample: Failed to listen rafthttp (%v)", err)
	}
	err = (&http.Server{Handler: rc.transport.Handler()}).Serve(ln)
	// todo: stopable
}

func NewRaftStorage(id uint64, peerIds []uint) {
	rn := &raftNode{}
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
	rn.node = raft.StartNode(c, peers)
	rn.transport = &rafthttp.Transport{
		ID:          types.ID(id),
		ClusterID:   0x1000,
		Raft:        rn,
		ServerStats: stats.NewServerStats("", ""),
		LeaderStats: stats.NewLeaderStats(strconv.FormatUint(id, 10)),
		ErrorC:      make(chan error),
	}
	// 配置transport: pipeline, stream
	rn.transport.Start()
	// raft http server, 节点间通讯用http?
	// /raft/stream -> stream: get -> peer请求Get就建立了，长连接
	// /raft -> pipeline: post raft message，stream不可用时
	// /raft/snapshot -> post snap message -> max(1<<63)
	// /raft/probing -> health(ok)
	go rn.serveRaft()
}
