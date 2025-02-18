package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WAQI Client
type WAQIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	cache   CacheService
}

func NewWAQIClient(apiKey string, baseURL string, cache CacheService) WAQIService {
	if cache == nil {
		cache = NewInMemoryCache()
	}

	return &WAQIClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		cache: cache,
	}
}

func (w *WAQIClient) GetByCity(ctx context.Context, city string) (*AirQuality, error) {
	cacheKey := GenerateAirQualityCacheKey(city)

	// Try to get from cache first
	if cached, found := w.cache.Get(cacheKey); found {
		if airQuality, ok := cached.(*AirQuality); ok {
			return airQuality, nil
		}
	}

	// If not in cache, fetch from API
	airQuality, err := w.fetchFromAPI(ctx, city)
	if err != nil {
		return nil, err
	}

	// Store in cache
	w.cache.Set(cacheKey, airQuality, DefaultCacheDuration)

	return airQuality, nil
}

func (w *WAQIClient) fetchFromAPI(ctx context.Context, city string) (*AirQuality, error) {
	uri := fmt.Sprintf("%s%s/?token=%s", w.baseURL, city, w.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s",
			resp.StatusCode, string(body))
	}

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result WaqiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Validate API response
	if result.Status != "ok" {
		return nil, fmt.Errorf("API returned non-ok status: %s", result.Status)
	}

	// Validate required data
	if len(result.Data.City.Geo) < 2 {
		return nil, fmt.Errorf("invalid geo data from API")
	}

	if len(result.Data.Attributions) == 0 {
		return nil, fmt.Errorf("missing station attribution data")
	}

	// Extract and format data
	airQuality := &AirQuality{
		City:     result.Data.City.Name,
		AQI:      result.Data.AQI,
		Lat:      result.Data.City.Geo[0],
		Lon:      result.Data.City.Geo[1],
		Station:  result.Data.Attributions[0].Name,
		UpdateAt: result.Data.Time.ISO,
	}

	// Validate constructed response
	if err := validateAirQuality(airQuality); err != nil {
		return nil, fmt.Errorf("invalid air quality data: %w", err)
	}

	return airQuality, nil
}

// validateAirQuality performs validation on the constructed AirQuality object
func validateAirQuality(aq *AirQuality) error {
	if aq.City == "" {
		return fmt.Errorf("city name is empty")
	}

	if aq.AQI < 0 {
		return fmt.Errorf("invalid AQI value: %d", aq.AQI)
	}

	if aq.Lat < -90 || aq.Lat > 90 {
		return fmt.Errorf("invalid latitude value: %f", aq.Lat)
	}

	if aq.Lon < -180 || aq.Lon > 180 {
		return fmt.Errorf("invalid longitude value: %f", aq.Lon)
	}

	if aq.Station == "" {
		return fmt.Errorf("station name is empty")
	}

	// Validate ISO timestamp format
	if _, err := time.Parse(time.RFC3339, aq.UpdateAt); err != nil {
		return fmt.Errorf("invalid timestamp format: %s", aq.UpdateAt)
	}

	return nil
}
