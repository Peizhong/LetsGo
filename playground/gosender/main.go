package main

import (
	"fmt"
	"log"
	"time"

	"github.com/peizhong/letsgo/playground/gosender/emitter"
)

const (
	connections    = 100
	milliseconds   = 100
	sizePerRequest = 512
	baseRoomId     = 1000
)

type message struct {
	RommId        string `msgpack:"roomId"`
	BroadcastData []int8 `msgpack:"broadcastData"`
	Time          int64  `msgpack:"time"`
}

func startup(count int, fn func(n int)) {
	started := 0
	step := 2
	threshold := count / 2
	for started < count {
		todo := count - started
		if todo > step {
			todo = step
		}
		for i := 0; i < todo; i++ {
			started++
			index := started
			go fn(index)
		}
		if started < threshold {
			step *= step
		}
		if started < count {
			<-time.After(time.Second)
		}
	}
}

func main() {
	opts := &emitter.Options{}
	socket := emitter.NewEmitter(opts)
	defer socket.Close()
	// 模拟发送数据，只用1个redis连接
	todo := make(chan string, connections)
	go func() {
		buf := make([]int8, sizePerRequest)
		msg := message{BroadcastData: buf, Time: time.Now().UnixNano() / 1000000}
		for roomId := range todo {
			socket.In(roomId).Broadcast().Emit("operator", msg)
		}
	}()
	startup(connections, func(index int) {
		log.Println("hello", index)
		roomId := fmt.Sprint(index)
		for range time.Tick(time.Millisecond * milliseconds) {
			todo <- roomId
		}
	})
	select {}
}
