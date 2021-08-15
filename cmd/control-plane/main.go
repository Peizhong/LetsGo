package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	server "github.com/envoyproxy/go-control-plane/pkg/server/v2"
	"github.com/peizhong/letsgo/cmd/control-plane/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

func registerServer(grpcServer *grpc.Server, server server.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("start serve stream", info.FullMethod)
		start := time.Now()
		err := handler(srv, stream)
		log.Println("stop steam, duration", time.Since(start).Seconds())
		return err
	}
}

// RunServer starts an xDS server at the given port.
func RunServer(ctx context.Context, srv server.Server, port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems. Keepalive timeouts based on connection_keepalive parameter https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples#dynamic
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.ChainStreamInterceptor(StreamServerInterceptor()),
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv)

	log.Printf("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}

func main() {
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, internal.NewLogger())
	ctx := context.Background()
	snapshot := internal.GenerateSnapshot()
	snapshotCache.SetSnapshot("", snapshot)
	srv := server.NewServer(ctx, snapshotCache, &internal.Callbacks{Debug: true})
	RunServer(ctx, srv, 18000)
}
