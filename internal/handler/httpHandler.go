package handler

import (
	"encoding/json"
	"net/http"

	"github.com/shii-park/Stshoot_Backend/internal/broadcaster"
)

func HandleCreate(w http.ResponseWriter, r *http.Request, hubManager *broadcaster.HubManager) {
	roomID, err := hubManager.CreateHub()
	if err != nil {
		http.Error(w, "Failed to create hub", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"roomID": roomID})
}
