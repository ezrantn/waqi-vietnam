# WAQI Vietnam

A simple REST API built with Go to fetch air quality data for major cities in Vietnam using the [World Air Quality Index (WAQI](https://waqi.info/#/c/5.59/7.129/2.7z) API.

## API Endpoint

### Air Quality Check

Fetches the air quality data for a given city in Vietnam.

**Request**

- URL: `/api/v1/air-quality/{city}`
- Method: GET
- Parameters:
  - city: The name of the city (case-insensitive) for which you want to fetch the air quality data.

**Response**

- Status Code: 200 OK (if successful)
- Response Body: A JSON object containing the air quality data for the requested city.

```json
{
    "status": "ok",
    "data": {
        "aqi": "-",
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
        }
    }
}
```

> [!IMPORTANT]
> Some cities may not provide AQI data directly. For example, Ho Chi Minh City’s AQI might be represented as "-". This typically occurs because the air quality monitoring station at that location is currently not providing any data.
> In such cases, the API may return an empty or default response for AQI-related values. It's recommended to check the status of air quality monitoring stations in those cities for potential data availability.

## Setup

To run the WAQI Vietnam API, follow the instructions below:

### Prerequisites

- Go: Install Go 1.23+.
- WAQI API Key: Sign up for a WAQI API key [here](https://aqicn.org/data-platform/token/).

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

3. Set your environment variable:

```bash
cp .env.example .env
```

Fill in the API token you received after signing up for the WAQI service and place it in the `API_TOKEN` field in the `.env` file.

4. Run the server:

```bash
go run .
```

### Example Usage

To get the air quality data for Hanoi:

```bash
curl http://localhost:3000/api/air-quality/hanoi
```

## License

This tool is open-source and available under the [MIT License](https://github.com/ezrantn/waqi-vietnam/blob/main/LICENSE).

## Contributing

Contributions are welcome! Please feel free to submit a pull request. For major changes, please open an issue first to discuss what you would like to change.