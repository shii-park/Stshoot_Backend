package main

import (
	"fmt"
	"net/http"

	"github.com/shii-park/Stshoot_Backend/internal/broadcaster"
	"github.com/shii-park/Stshoot_Backend/internal/handler"
)

func main() {

	receiver := &broadcaster.Receiver{}

	hub := broadcaster.NewHub(receiver)
	go hub.Run()

	http.HandleFunc("/ws/receiver", func(w http.ResponseWriter, r *http.Request) { handler.HandleReceiver(w, r, receiver) })
	http.HandleFunc("/ws/sender", func(w http.ResponseWriter, r *http.Request) { handler.HandleSender(w, r, hub) })
	fmt.Println("WebSocket server started on ws://localhost:80/ws")
	http.ListenAndServe(":80", nil)
}
