package broadcaster

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Broadcast    chan *hubMessage
	Register     chan *SenderClient
	Unregister   chan *SenderClient
	Clients      map[*SenderClient]bool
	stop         chan struct{}
	Receiver     *Receiver
	lastActivity time.Time
	mu           sync.RWMutex
}

func NewHub(receiver *Receiver) *Hub {
	return &Hub{
		Broadcast:    make(chan *hubMessage),
		Register:     make(chan *SenderClient),
		Unregister:   make(chan *SenderClient),
		Clients:      make(map[*SenderClient]bool),
		stop:         make(chan struct{}),
		Receiver:     receiver,
		lastActivity: time.Now(),
	}
}

type hubMessage struct {
	data   []byte
	sender *SenderClient
}

// Hubのメインループ
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register: // Senderがコネクションを行ったとき
			h.Clients[client] = true
			log.Printf("Sender client connected. Total senders: %d", len(h.Clients))
		case client := <-h.Unregister: // Senderがコネクションから切断されたとき
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				log.Printf("Sender client disconnected. Total senders: %d", len(h.Clients))
			}
		case message := <-h.Broadcast: // メッセージが送られてきたとき
			log.Printf("broadcast message: %s", message.data)
			for client := range h.Clients {
				if client != message.sender {
					client.Send <- message.data
				}
			}
			if err := h.Receiver.send(message.data); err != nil {
				log.Printf("Could not send message to receiver: %v", err)
			}
		case <-h.stop: // stopチャネルが閉じられたとき
			log.Printf("Hub for receiver %p is stopping.", h.Receiver)
			for client := range h.Clients {
				close(client.Send)
			}
			return // forループを抜けてゴルーチンを終了する
		}
	}
}

// Receiverに対してJSONを送信する
func (r *Receiver) send(message []byte) error {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if r.Conn == nil {
		return fmt.Errorf("no receiver client connected")
	}
	return r.Conn.WriteMessage(websocket.TextMessage, message)
}

// Hubに人がいないかどうかを確かめる
func (h *Hub) IsEmpty() bool {
	h.Receiver.Mu.RLock() // Receiverの状態を安全に読む
	defer h.Receiver.Mu.RUnlock()

	h.mu.RLock() // Hubの状態を安全に読む
	defer h.mu.RUnlock()

	return len(h.Clients) == 0 && h.Receiver.Conn == nil
}
