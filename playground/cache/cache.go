package main

import (
	"github.com/go-redis/redis"
	etcd "go.etcd.io/etcd/client"

	"context"
	"log"
	"time"
)

type Cacher interface {
	Init() error
	GetConfig() map[string]interface{}

	SetString(key, value string) error
	GetString(key string) (value string, err error)
	// Del
}

type GoRedis struct {
	client *redis.Ring
}

func (r *GoRedis) Init() error {
	// 单节点
	/*client := redis.NewClient(&redis.Options{
		Addr: "193.112.41.28:6379",
		DB:   0,
	})
	*/
	// 分区
	client := redis.NewRing(&redis.RingOptions{
		Addrs:    map[string]string{"main": "193.112.41.28:6379"},
		DB:       0,
		Password: "ur@hello123",
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Panicln(err)
	}
	log.Println(pong)
	r.client = client
	return nil
}

func (r *GoRedis) GetConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (r *GoRedis) SetString(key, value string) error {
	err := r.client.Set(key, value, time.Second).Err()
	return err
}

func (r *GoRedis) GetString(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		// log.Println("key does not exist")
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

type GoEtcd struct {
	client *etcd.Client
	kapi   etcd.KeysAPI
}

func (e *GoEtcd) Init() error {
	cfg := etcd.Config{
		Endpoints: []string{"http://193.112.41.28:2379"},
		Transport: etcd.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		return err
	}
	kapi := etcd.NewKeysAPI(client)
	e.client = &client
	e.kapi = kapi

	return nil
}

func (e *GoEtcd) GetConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (e *GoEtcd) SetString(key, value string) error {
	_, err := e.kapi.Set(context.Background(), key, value, &etcd.SetOptions{TTL: time.Second})
	if err != nil {
		return err
	}
	// print common key info
	// log.Printf("Set is done. Metadata is %q\n", resp)
	return nil
}

func (e *GoEtcd) GetString(key string) (string, error) {
	resp, err := e.kapi.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	}
	// log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	return resp.Node.Value, nil
}
