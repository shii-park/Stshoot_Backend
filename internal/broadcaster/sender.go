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
	Send chan []byte
}

func (c *SenderClient) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		}
		log.Printf("c.sent message: %s", message)
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error writing message: %v", err)
			return
		}
	}
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
		log.Printf("sent message %s", msg.Text)
		if err != nil {
			break
		}

		jsonBytes, err := json.Marshal(msg)

		if err != nil {
			log.Printf("error marshalling json: %v", err)
			continue
		}

		hubMessage := &hubMessage{data: jsonBytes, sender: c}

		c.Hub.Broadcast <- hubMessage
	}
}
