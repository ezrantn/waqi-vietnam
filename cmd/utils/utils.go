package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	BASE_URL  string
	API_TOKEN string
)

var VietnamCities = []string{
	"hanoi",
	"ho-chi-minh",
	"da-nang",
	"haiphong",
	"can-tho",
	"nha-trang",
	"dalat",
	"hue",
	"vinh",
	"quy-nhon",
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	BASE_URL = os.Getenv("BASE_URL")
	API_TOKEN = os.Getenv("API_TOKEN")

	if BASE_URL == "" || API_TOKEN == "" {
		log.Println("Warning: Required environment variables are not set")
	}
}
