package broadcaster

import "log"

type Hub struct {
	Clients    map[*SenderClient]bool
	Broadcast  chan []byte
	Register   chan *SenderClient
	Unregister chan *SenderClient
	Receiver   *Receiver
}

func NewHub(receiver *Receiver) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *SenderClient),
		Unregister: make(chan *SenderClient),
		Clients:    make(map[*SenderClient]bool),
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
		}
	}
}
