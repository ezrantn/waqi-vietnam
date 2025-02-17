package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type Utils struct{}

var (
	BASE_URL  string
	API_TOKEN string
	PORT      string
)

// We apply a rate limiter to prevent abuse of the /api/v1/air-quality/ endpoint.
// Although the third-party API allows up to 1000 requests per second, we narrow it down to:
// - Allow 1 request per second, with a burst capacity of 5 requests.
// This helps ensure fair usage, protects our infrastructure from traffic spikes,
// reduces unnecessary API calls, and improves overall system stability.
var limiter = rate.NewLimiter(1, 5)

func (u *Utils) CorsMiddleware(next http.Handler) http.Handler {
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

func (u *Utils) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var VietnamCities = []string{
	"hanoi",
	"ho-chi-minh",
	"da-nang",
	"haiphong",
	"can-tho",
	"nha-trang",
	"hue",
	"vinh",
	"thai-nguyen",
}

func (u *Utils) NormalizeCity(city string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(city), " ", "-"))
}

func (u *Utils) IsValidVietnamCity(city string) bool {
	normalizedCity := u.NormalizeCity(city)
	for _, v := range VietnamCities {
		if normalizedCity == v {
			return true
		}
	}
	return false
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	BASE_URL = os.Getenv("BASE_URL")
	API_TOKEN = os.Getenv("API_TOKEN")
	PORT = os.Getenv("PORT")

	if BASE_URL == "" || API_TOKEN == "" || PORT == "" {
		log.Println("Warning: Required environment variables are not set")
	}
}
