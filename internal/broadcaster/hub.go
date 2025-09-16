package broadcaster

import "log"

type Hub struct {
	Broadcast  chan []byte
	Register   chan *SenderClient
	Unregister chan *SenderClient
	Clients    map[*SenderClient]bool
	stop       chan struct{}
	Receiver   *Receiver
}

func NewHub(receiver *Receiver) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *SenderClient),
		Unregister: make(chan *SenderClient),
		Clients:    make(map[*SenderClient]bool),
		stop:       make(chan struct{}),
		Receiver:   receiver,
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
