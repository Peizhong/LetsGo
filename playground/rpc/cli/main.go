package main

import (
	"context"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/helloworld"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := pb.NewGreeterClient(conn)
	r, err := cli.SayHello(context.Background(), &pb.HelloRequest{
		Name: "client",
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(r)
}
