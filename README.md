# WAQI Vietnam API - Backend

A robust REST API built with Go that provides real-time air quality data for major cities in Vietnam through the [World Air Quality Index (WAQI](https://waqi.info/#/c/5.59/7.129/2.7z) API. This service powers the WAQI Research project, combining Go's powerful backend capabilities with React's dynamic frontend to deliver comprehensive air quality monitoring and analysis.

## Features

- Real-time air quality data retrieval
- Support for all major Vietnamese cities
- JSON-based REST API
- Efficient data caching
- Rate limiting protection
- Comprehensive error handling
- CORS support for web applications

## API Reference

### Get Air Quality Data

Retrieves current air quality data for a specified Vietnamese city.

### Endpoint

```bash
GET /api/v1/air-quality/{city}
```

### Parameters

| Parameter | Type   | Required | Description                                 |
|-----------|--------|----------|---------------------------------------------|
| city      | string | yes      | City name (case-insensitive, e.g., "hanoi") |

### Response

**Success Response (200 OK)**

```json
{
    "status": "ok",
    "data": {
        "aqi": "75",
        "attributions": [
            {
                "url": "https://vn.usembassy.gov/embassy-consulates/ho-chi-minh-city/air-quality-monitor/",
                "name": "Ho Chi Minh City Air Quality Monitor - Embassy of the United States"
            },
            {
                "url": "https://waqi.info/",
                "name": "World Air Quality Index Project"
            }
        ],
        "city": {
            "name": "Ho Chi Minh City US Consulate, Vietnam (Lãnh sự quán Hoa Kỳ, Hồ Chí Minh)",
            "geo": [
                10.782978,
                106.700711
            ]
        },
        "timestamp": "2024-02-17T08:00:00Z",
        "dominentPollutant": "PM2.5",
        "iaqi": {
            "pm25": {
                "v": 75
            },
            "humidity": {
                "v": 65
            },
            "temperature": {
                "v": 28
            }
        }
    }
}
```

### Error Responses

- 400 Bad Request: Invalid city name
- 404 Not Found: City not found
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server-side error
- 503 Service Unavailable: WAQI service unavailable

### Data Interpolation

The API returns several key metrics:

- `aqi`: Air Quality Index value (0-500)
- `dominentPollutant`: Main pollutant affecting air quality
- `iaqi`: Individual air quality parameters
  - `pm25`: Fine particulate matter
  - `humidity`: Relative humidity
  - `temperature`: Ambient temperature

> [!IMPORTANT]
> Some monitoring stations may temporarily report missing data, indicated by "-" in the AQI field. This usually means the station is undergoing maintenance or experiencing technical issues. Applications should handle these cases gracefully by displaying appropriate user messages.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- WAQI API token (obtain from [WAQI Data Platform](https://aqicn.org/data-platform/token/))
- Git

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ezrantn/waqi-vietnam.git
   cd waqi-vietnam
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure environment:

   ```bash
   cp .env.example .env
   ```

4. Add your WAQI API token to `.env`:

   ```env
   BASE_URL="https://api.waqi.info/feed/"
   API_TOKEN=your_token_here
   PORT="3000"
   ```

5. Start the server:

   ```bash
   go run .
   ```

### Example Usage

To get the air quality data for Hanoi:

```bash
curl http://localhost:3000/api/v1/air-quality/hanoi
```

## License

This tool is open-source and available under the [MIT License](https://github.com/ezrantn/waqi-vietnam/blob/main/LICENSE).

## Contributing

Contributions are welcome! Please feel free to submit a pull request. For major changes, please open an issue first to discuss what you would like to change.
