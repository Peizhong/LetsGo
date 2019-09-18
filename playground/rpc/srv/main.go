package main

import (
	"context"
	"github.com/peizhong/letsgo/pkg/config"
	"github.com/peizhong/letsgo/playground/rpc/pb/helloworld"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/twoway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
)

type server struct {
}

func (s *server) DoBothWay(context.Context, *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello"}, nil
}

func (s *server) Say2Hello(context.Context, *pb.Hello2Request) (*pb.Hello2Reply, error) {
	return &pb.Hello2Reply{Message: "hello"}, nil
}

func main() {
	port := config.GrpcApp1Port
	addr := net.JoinHostPort("", strconv.Itoa(port))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listening at", addr)
	creds, err := credentials.NewServerTLSFromFile(config.CertCrt, config.CertKey)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterTwoWayJobServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
