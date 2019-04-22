package catalog

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/peizhong/letsgo/framework"
	pb "github.com/peizhong/letsgo/services/helloworld"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func hiGrpc() {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case runtime.Error:
				log.Printf("runtime error: %v", r)
			default:
				log.Printf("error: %v", r)
			}
		}
	}()
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

func GetProduct(c *framework.GoContext, r GetProductRequest) (*GetProductResponse, error) {
	product := new(Product)
	db, err := c.GetDatabase()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer db.Close()
	db.Where(&Product{Id: r.ProductId}).First(product)

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

	response := new(GetProductResponse)
	framework.DirectMapTo(product, response)
	return response, nil
}

func GetProducts(c *framework.GoContext, r GetProductsRequest) (*GetProductsResponse, error) {
	var products []*Product
	db, err := c.GetDatabase()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer db.Close()

	// raw sql
	db.Raw("SELECT SQL_CALC_FOUND_ROWS * FROM `Products`;").Scan(&products)
	// count := make([]int, 0)
	// db.Raw("SELECT FOUND_ROWS() ROWCOUNT;").Pluck("ROWCOUNT", &count)
	var count int
	db.Raw("SELECT FOUND_ROWS() ROWCOUNT;").Row().Scan(&count)

	db.Offset(r.PageIndex * r.PageSize).Limit(r.PageSize).Find(&products)
	response := &GetProductsResponse{
		Count:     count,
		PageIndex: r.PageIndex,
		PageSize:  r.PageSize,
	}
	tmp := struct {
		Items []*Product
	}{Items: products}
	framework.DirectMapTo(&tmp, response)
	return response, nil
}
