package main

import (
	"fmt"
	"net/http"

	"github.com/shii-park/Stshoot_Backend/internal/handler"
)

func main() {
	http.HandleFunc("/ws", handler.HandleWebSocket)
	fmt.Println("WebSocket server started on ws://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}
