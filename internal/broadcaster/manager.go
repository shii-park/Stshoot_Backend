package broadcaster

import (
	"fmt"
	"log"
	"sync"

	"github.com/shii-park/Stshoot_Backend/internal/utils"
)

type HubManager struct {
	Hubs map[string]*Hub
	Mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	m := &HubManager{
		Hubs: make(map[string]*Hub),
	}
	go m.runJanitor()
	return m
}

// roomIDから*Hubを取得する
func (m *HubManager) GetHub(roomID string) (*Hub, error) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	if hub, ok := m.Hubs[roomID]; ok {
		return hub, nil
	}
	return nil, fmt.Errorf("hub with id '%s' not found", roomID)

}

// Hubを作成する
func (m *HubManager) CreateHub() (string, error) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	for {
		roomID, err := utils.GenRandomID(6)
		if err != nil {
			return "", err
		}
		if _, ok := m.Hubs[roomID]; !ok {
			receiver := &Receiver{}
			hub := NewHub(receiver)
			go hub.Run()
			m.Hubs[roomID] = hub // マップに新しいHubを登録
			log.Printf("New hub created for room: %s", roomID)
			return roomID, nil
		}
	}

}

func (m *HubManager) DeleteHub(roomID string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	if hub, ok := m.Hubs[roomID]; ok {
		close(hub.stop)

		delete(m.Hubs, roomID)
		log.Printf("Hub for room '%s' has been deleted.", roomID)
	}
}
