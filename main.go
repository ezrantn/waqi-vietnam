package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ezrantn/waqivietnam/cmd/api"
)

func main() {
	http.HandleFunc("/api/air-quality/", api.AirQualityHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
