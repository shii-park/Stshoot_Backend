package broadcaster

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type SenderClient struct {
	Hub  *Hub
	Conn *websocket.Conn
}

func (r *Receiver) Send(message []byte) error {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if r.Conn == nil {
		return fmt.Errorf("no receiver client connected")
	}
	return r.Conn.WriteMessage(websocket.TextMessage, message)
}

func (c *SenderClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// メッセージに送信元情報を付加する
		senderInfo := fmt.Sprintf("[%s says]: ", c.Conn.RemoteAddr())
		fullMessage := append([]byte(senderInfo), message...)
		c.Hub.Broadcast <- fullMessage
	}
}
