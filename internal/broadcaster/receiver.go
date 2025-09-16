package broadcaster

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Receiver struct {
	Conn *websocket.Conn
	Mu   sync.RWMutex
}

func (r *Receiver) SetConnection(conn *websocket.Conn) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if r.Conn != nil {
		return fmt.Errorf("a receiver client is already connected")
	}
	r.Conn = conn
	log.Println("Receiver client connected.")
	return nil
}

func (r *Receiver) ClearConnection() {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if r.Conn != nil {
		r.Conn.Close()
		r.Conn = nil
	}
	log.Println("Receiver client disconnected.")
}
