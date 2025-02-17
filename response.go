package main

// WAQI API Response
type WaqiResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI  int `json:"aqi"`
		City struct {
			Name string    `json:"name"`
			Geo  []float64 `json:"geo"`
		} `json:"city"`
		Idx  int `json:"idx"`
		Time struct {
			ISO string `json:"iso"`
		} `json:"time"`
		Forecast struct {
			Daily map[string][]struct {
				Avg int `json:"avg"`
			} `json:"daily"`
		} `json:"forecast"`
		Attributions []struct {
			Name string `json:"name"`
		} `json:"attributions"`
	} `json:"data"`
}
