package main

import (
	"context"
	"github.com/go-redis/redis"
	"time"

	etcd "go.etcd.io/etcd/client"
	"log"
)

func redisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "193.112.41.28:6379",
		Password: "ur@hello123",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Panicln(err)
	}
	log.Println(pong)
	err = client.Set("key", "value", 0).Err()
	if err != nil {
		log.Println(err)
	}
	val, err := client.Get("key").Result()
	if err == redis.Nil {
		log.Println("key does not exist")
	} else if err != nil {
		log.Panicln(err)
	} else {
		log.Println(val)
	}
}

func etcdClient() {
	cfg := etcd.Config{
		Endpoints: []string{"http://193.112.41.28:2379"},
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		log.Panicln(err)
	}
	kapi := etcd.NewKeysAPI(client)
	resp, err := kapi.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		log.Panicln(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
	resp, err = kapi.Get(context.Background(), "/foo", nil)
	if err != nil {
		log.Panicln(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
}

func main() {
	redisClient()
	etcdClient()
}
