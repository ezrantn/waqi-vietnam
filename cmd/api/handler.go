package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ezrantn/waqivietnam/cmd/utils"
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

func FetchAirQuality(city string) (*WaqiResponse, error) {
	if utils.API_TOKEN == "" || utils.BASE_URL == "" {
		return nil, fmt.Errorf("API_TOKEN or BASE_URL is not set")
	}

	apiToken := utils.API_TOKEN
	if apiToken == "" {
		return nil, fmt.Errorf("WAQI TOKEN is not set")
	}

	uri := fmt.Sprintf("%s%s/?token=%s", utils.BASE_URL, city, apiToken)

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

// Handler function for API requests
func AirQualityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract city name from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/air-quality/")
	if path == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	// Validate city name
	valid := false
	for _, v := range utils.VietnamCities {
		if v == path {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, "Invalid city", http.StatusBadRequest)
		return
	}

	// Fetch air quality data
	data, err := FetchAirQuality(path)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Set JSON response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	json.NewEncoder(w).Encode(data)
}
