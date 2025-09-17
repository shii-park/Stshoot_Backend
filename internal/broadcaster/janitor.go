package broadcaster

import (
	"log"
	"time"
)

func (m *HubManager) runJanitor() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	var inactiveHubIDs []string

	for range ticker.C {

		m.Mu.RLock()
		for id, hub := range m.Hubs {
			hub.mu.RLock()

			if hub.IsEmpty() && time.Since(hub.lastActivity) > 5*time.Minute {
				inactiveHubIDs = append(inactiveHubIDs, id)
			}
			hub.mu.RUnlock()
		}
		m.Mu.RUnlock()
	}
	for _, id := range inactiveHubIDs {
		log.Printf("Janitor: Deleting inactive hub for room '%s'", id)
		m.DeleteHub(id)
	}

}
