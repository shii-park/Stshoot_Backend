package broadcaster

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	Broadcast    chan []byte
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
		Broadcast:    make(chan []byte),
		Register:     make(chan *SenderClient),
		Unregister:   make(chan *SenderClient),
		Clients:      make(map[*SenderClient]bool),
		stop:         make(chan struct{}),
		Receiver:     receiver,
		lastActivity: time.Now(),
	}
}

// Hubのメインループ
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Sender client connected. Total senders: %d", len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				log.Printf("Sender client disconnected. Total senders: %d", len(h.Clients))
			}
		case message := <-h.Broadcast:
			if err := h.Receiver.Send(message); err != nil {
				log.Printf("Could not send message to receiver: %v", err)
			}
		case <-h.stop: // stopチャネルが閉じられたら
			log.Printf("Hub for receiver %p is stopping.", h.Receiver)
			return // forループを抜けてゴルーチンを終了する
		}
	}
}

func (h *Hub) IsEmpty() bool {
	h.Receiver.Mu.RLock() // Receiverの状態を安全に読む
	defer h.Receiver.Mu.RUnlock()

	h.mu.RLock() // Hubの状態を安全に読む
	defer h.mu.RUnlock()

	return len(h.Clients) == 0 && h.Receiver.Conn == nil
}
