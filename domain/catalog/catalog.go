package catalog

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/peizhong/letsgo/framework"
	pb "github.com/peizhong/letsgo/services/helloworld"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func hiGrpc() {
	const (
		address     = "localhost:50051"
		defaultName = "world"
	)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	demo := []int{1, 2}
	for d := range demo {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("%v_%v", d, name)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting %d: %s", d, r.Message)
		r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: fmt.Sprintf("%v_%v", d, name)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting again %d: %s", d, r.Message)
	}
}

func GetProduct(c *framework.GoContext, r GetProductRequest) (*GetProductsResponse, error) {
	var classify Classify
	db, err := c.GetDatabase()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer db.Close()
	db.First(&classify, r.ProductId)
	// play redis
	cache, _ := c.GetCache()
	defer cache.Close()
	_, err = cache.Do("SET", "go_key", "redigo")
	if err == nil {
		v, err := redis.String(cache.Do("GET", "go_key"))
		if err == nil {
			_ = v
		}
	}
	// play grpc
	hiGrpc()
	return &GetProductsResponse{
		PageIndex: 1,
		PageSize:  100,
	}, nil
}

func GetProducts(c *framework.GoContext, r GetProductsRequest) (*GetProductsResponse, error) {
	var classifies []Classify
	db, err := c.GetDatabase()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer db.Close()
	db.Offset(r.PageIndex * r.PageSize).Limit(r.PageSize).Find(&classifies)
	return &GetProductsResponse{
		PageIndex: r.PageIndex,
		PageSize:  r.PageSize,
	}, nil
}
