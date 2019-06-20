package main

import (
	"github.com/peizhong/letsgo/internal"
	"github.com/tidwall/evio"
)

var events evio.Events

func Start() {
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		res := string(in)
		res += "woo\r\n"
		out = []byte(res)
		return
	}
	if err := evio.Serve(events, "tcp://localhost:8082"); err != nil {
		panic(err.Error())
	}
}

func main() {
	internal.Host(Start, func() {

	})
}
