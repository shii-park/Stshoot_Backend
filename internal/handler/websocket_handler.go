package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket connection failed:", err)
		return
	}

	go handleConnection(conn) //一つのコネクションに対してゴルーチンで並列にあつかう
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	welcomeMessage := []byte("サーバーへの接続に成功しました！")
	println("Someone connected")

	if err := conn.WriteMessage(websocket.TextMessage, welcomeMessage); err != nil {

		fmt.Println("Failed to write welcome message:", err)
		return
	}
}
