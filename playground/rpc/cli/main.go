package main

import (
	"context"
	"github.com/peizhong/letsgo/pkg/config"
	pb "github.com/peizhong/letsgo/playground/rpc/pb/twoway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
)

func main() {
	port := config.GrpcApp1Port
	addr := net.JoinHostPort("", strconv.Itoa(port))
	creds, _ := credentials.NewClientTLSFromFile(config.CertCrt, config.CertName)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	cli := pb.NewTwoWayJobClient(conn)
	r, err := cli.Say2Hello(context.Background(), &pb.Hello2Request{
		Name: "client",
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(r)
}
