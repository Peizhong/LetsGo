package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

// redis 的订阅发布，一般用rabbit，看看客户端实现方式

func pub() {
	redis := redis.NewClient(&redis.Options{
		Addr:     "193.112.41.28:6379",
		Password: "ur@hello123",
		DB:       0,
	})
	fmt.Println("redis publish client:")
	for {
		var data string
		fmt.Scanln(&data)
		err := redis.Publish("mychannel1", data).Err()
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("ok")
	}
}

func main() {
	pub()
}
