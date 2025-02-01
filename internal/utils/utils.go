package utils

import "strings"

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

func NormalizeCity(city string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(city), " ", "-"))
}

func IsValidVietnamCity(city string) bool {
	normalizedCity := NormalizeCity(city)
	for _, v := range VietnamCities {
		if normalizedCity == v {
			return true
		}
	}
	return false
}
