package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// HTTP Handlers
type Handler struct {
	waqiClient *WAQIClient
	utils      *Utils
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *Handler) GetAirQualityByCity(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimPrefix(r.URL.Path, "/api/v1/air-quality/")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	normalizedCity := h.utils.NormalizeCity(city)
	if !h.utils.IsValidVietnamCity(normalizedCity) {
		http.Error(w, "Invalid city", http.StatusBadRequest)
		return
	}

	data, err := h.waqiClient.GetByCity(city)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Discussion(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	defer conn.Close()

	// Listen for incoming message
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
