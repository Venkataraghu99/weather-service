# Weather Service

A simple HTTP server that provides weather forecasts using the National Weather Service API.

## Features

- Get weather forecast by latitude and longitude coordinates
- Returns short forecast description
- Provides temperature condition (hot, cold, or moderate)
- Uses the official National Weather Service API

## Prerequisites

- Go 1.21 or higher

## Installation

1. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

```bash
go run cmd/weather-service/main.go
```

The server will start on `http://localhost:8080`.

## API Endpoint

### Get Weather Forecast

**Endpoint:** `GET /api/weather/forecast`

**Query Parameters:**
- `lat`: Latitude (required)
- `long`: Longitude (required)

**Example Request:**
```
GET /api/weather/forecast?lat=40.7128&long=-74.0060
```

**Example Response:**
```json
{
  "short_forecast": "Mostly Sunny",
  "temperature": 72,
  "unit": "F",
  "condition": "moderate"
}
```

**Error Responses:**
- `400 Bad Request` - Missing or invalid parameters
- `500 Internal Server Error` - Error fetching weather data

## Temperature Conditions

- **Hot**: 85°F (29.4°C) or above
- **Cold**: 45°F (7.2°C) or below
- **Moderate**: Between 46°F (7.8°C) and 84°F (28.9°C)


