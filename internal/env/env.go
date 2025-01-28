package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BASE_URL  string
	API_TOKEN string
	PORT      string
)

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
