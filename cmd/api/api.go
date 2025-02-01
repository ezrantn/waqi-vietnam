package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ezrantn/waqivietnam/internal/env"
	"github.com/ezrantn/waqivietnam/internal/utils"
)

type WaqiResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI          any `json:"aqi"`
		Attributions []struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"attributions"`
		City struct {
			Name string    `json:"name"`
			Geo  []float64 `json:"geo"`
		} `json:"city"`
	} `json:"data"`
}

// Helper function for fetching air quality
func fetchAirQuality(city string) (*WaqiResponse, error) {
	if env.API_TOKEN == "" || env.BASE_URL == "" {
		return nil, fmt.Errorf("API_TOKEN or BASE_URL is not set")
	}

	apiToken := env.API_TOKEN
	if apiToken == "" {
		return nil, fmt.Errorf("WAQI TOKEN is not set")
	}

	uri := fmt.Sprintf("%s%s/?token=%s", env.BASE_URL, city, apiToken)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WaqiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("API response status is: %s", result.Status)
	}

	return &result, nil
}

func AirQualityHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimPrefix(r.URL.Path, "/api/air-quality/")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	normalizedCity := utils.NormalizeCity(city)
	if !utils.IsValidVietnamCity(normalizedCity) {
		http.Error(w, "Invalid city", http.StatusBadRequest)
		return
	}

	data, err := fetchAirQuality(city)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
