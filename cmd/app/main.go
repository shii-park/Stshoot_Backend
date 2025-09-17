package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/shii-park/Stshoot_Backend/internal/broadcaster"
	"github.com/shii-park/Stshoot_Backend/internal/handler"
)

func main() {

	hubManager := broadcaster.NewHubManager()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws/sender/", func(w http.ResponseWriter, r *http.Request) {
		roomID := strings.TrimPrefix(r.URL.Path, "/ws/sender/")
		if roomID == "" {
			http.Error(w, "Room ID is required", http.StatusBadRequest)
			return
		}

		hub, err := hubManager.GetHub(roomID)
		if err != nil {
			http.Error(w, "Hub not found: "+err.Error(), http.StatusNotFound)
			return
		}
		handler.HandleSender(w, r, hub)
	})
	mux.HandleFunc("/ws/receiver/", func(w http.ResponseWriter, r *http.Request) {
		roomID := strings.TrimPrefix(r.URL.Path, "/ws/receiver/")
		if roomID == "" {
			http.Error(w, "Room ID is required", http.StatusBadRequest)
			return
		}
		hub, err := hubManager.GetHub(roomID)
		if err != nil {
			http.Error(w, "Hub not found: "+err.Error(), http.StatusNotFound)
			return
		}
		receiver := hub.Receiver
		handler.HandleReceiver(w, r, receiver, hubManager, roomID)
	})
	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleCreate(w, r, hubManager)
	})
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Server is active")

	})
	fmt.Println("WebSocket server started")
	if err := http.ListenAndServe(":10000", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
