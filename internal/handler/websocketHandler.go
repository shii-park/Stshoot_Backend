package handler

import (
	"log"
	"net/http"

	"github.com/shii-park/Stshoot_Backend/internal/broadcaster"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleSender(w http.ResponseWriter, r *http.Request, hub *broadcaster.Hub) {
	conn, err := upgrader.Upgrade(w, r, nil) // wsにアップグレード
	if err != nil {
		log.Printf("Failed to upgrade sender connection: %v", err)
		return
	}
	client := &broadcaster.SenderClient{Hub: hub, Conn: conn}
	client.Hub.Register <- client

	go client.ReadPump()
}

func HandleReceiver(w http.ResponseWriter, r *http.Request, receiver *broadcaster.Receiver) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade receiver connection: %v", err)
		return
	}
	if err := receiver.SetConnection(conn); err != nil {
		log.Printf("Failed to set receiver connection: %v", err)
		conn.Close()
		return
	}
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			receiver.ClearConnection()
			break
		}
	}
}
