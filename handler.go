package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
