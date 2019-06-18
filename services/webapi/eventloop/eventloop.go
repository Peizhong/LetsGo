package eventloop

import "github.com/tidwall/evio"

type EventLoop struct {
}

var events evio.Events

func (*EventLoop) Start() {
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

func (*EventLoop) Stop() {

}
