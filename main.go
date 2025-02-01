package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ezrantn/waqivietnam/cmd/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/air-quality/", api.AirQualityHandler)
	mux.HandleFunc("/api/health", api.HealthCheckHandler)

	handler := api.CorsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
