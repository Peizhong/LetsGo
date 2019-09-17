package main

import (
	"context"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (s *server) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello"}, nil
}

func main() {
	lis, err := net.Listen("tcp", "8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
