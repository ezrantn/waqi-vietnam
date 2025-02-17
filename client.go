package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// WAQI Client
type WAQIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewWAQIClient(apiKey string) *WAQIClient {
	return &WAQIClient{
		apiKey:  apiKey,
		baseURL: os.Getenv("BASE_URL"),
		client:  &http.Client{},
	}
}

func (w *WAQIClient) GetByCity(city string) (*AirQuality, error) {
	if w.apiKey == "" || w.baseURL == "" {
		return nil, fmt.Errorf("API_KEY or BASE_URL is not set")
	}

	uri := fmt.Sprintf("%s%s/?token=%s", w.baseURL, city, w.apiKey)

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

	// Check if geo has valid data
	if len(result.Data.City.Geo) < 2 {
		return nil, fmt.Errorf("invalid geo data from API")
	}

	// Extract relevant data and return
	airQuality := &AirQuality{
		City:     result.Data.City.Name,
		AQI:      result.Data.AQI,
		Lat:      result.Data.City.Geo[0], // First element is latitude
		Lon:      result.Data.City.Geo[1], // Second element is longitude
		Station:  result.Data.Attributions[0].Name,
		UpdateAt: result.Data.Time.ISO,
	}

	return airQuality, nil
}
