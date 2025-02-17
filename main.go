package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("API_TOKEN")
	if apiKey == "" {
		log.Fatal("API_TOKEN environment variable is required")
	}

	waqiClient := NewWAQIClient(apiKey)
	handler := &Handler{waqiClient: waqiClient}
	u := &Utils{}

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", handler.HealthCheck)
	mux.Handle("/api/v1/air-quality/", u.CorsMiddleware(u.RateLimit(http.HandlerFunc(handler.GetAirQualityByCity))))
	mux.Handle("/api/v1/discussion", u.CorsMiddleware(u.RateLimit(http.HandlerFunc(handler.Discussion))))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
