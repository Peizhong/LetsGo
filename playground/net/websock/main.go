package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func echo(w *websocket.Conn) {
	var reply string
	for {
		if err := websocket.Message.Receive(w, &reply); err != nil {
			log.Printf("receive error: %s\n", err.Error())
			return
		}
		log.Printf("message:%s\n", reply)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	// var ws = new WebSocket("ws://localhost:8081/v2/ws")
	// ws.addEventListener("message",function(e){console.log(e);});
	// ws.send("hello")
	http.Handle("/v2/ws", websocket.Handler(echo))
	http.ListenAndServe(":8081", nil)
}
