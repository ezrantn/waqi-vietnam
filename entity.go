package main

// Entity
type AirQuality struct {
	City     string  `json:"city"`
	AQI      int     `json:"aqi"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Station  string  `json:"station"`
	UpdateAt string  `json:"update_at"`
}
