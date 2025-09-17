package broadcaster

import (
	"encoding/json"
	"log"

	"github.com/shii-park/Stshoot_Backend/internal/model"

	"github.com/gorilla/websocket"
)

type SenderClient struct {
	Hub  *Hub
	Conn *websocket.Conn
}

// メッセージ読み取り
func (c *SenderClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg model.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		jsonBytes, err := json.Marshal(msg)

		if err != nil {
			log.Printf("error marshalling json: %v", err)
			continue
		}
		c.Hub.Broadcast <- jsonBytes
	}
}
