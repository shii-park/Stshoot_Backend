package broadcaster

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shii-park/Stshoot_Backend/internal/model"

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
