package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	client *http.Client
}

func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number          int    `json:"number"`
			Name            string `json:"name"`
			ShortForecast   string `json:"shortForecast"`
			Temperature     int    `json:"temperature"`
			TemperatureUnit string `json:"temperatureUnit"`
			IsDaytime       bool   `json:"isDaytime"`
			StartTime       string `json:"startTime"`
		} `json:"periods"`
	} `json:"properties"`
}

type ForecastResult struct {
	ShortForecast string `json:"short_forecast"`
	Temperature   int    `json:"temperature"`
	Unit          string `json:"unit"`
	Condition     string `json:"condition"`
}

func (s *Service) GetForecast(lat, long string) (*ForecastResult, error) {
	// First, get the forecast endpoint for the given coordinates
	pointURL := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", parseFloat(lat), parseFloat(long))

	headers := http.Header{
		"User-Agent": {"weather-service/1.0 (your-email@example.com)"},
		"Accept":     {"application/geo+json"},
	}

	// Get the forecast URL for the point
	forecastURL, err := s.getForecastURL(pointURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast URL: %w", err)
	}

	// Get the forecast data
	req, err := http.NewRequest("GET", forecastURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = headers

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body)
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("failed to decode forecast: %w", err)
	}

	// Find today's forecast (first period where isDaytime is true)
	for _, period := range forecast.Properties.Periods {
		if period.IsDaytime {
			return &ForecastResult{
				ShortForecast: period.ShortForecast,
				Temperature:   period.Temperature,
				Unit:          period.TemperatureUnit,
				Condition:     getTemperatureCondition(period.Temperature, period.TemperatureUnit),
			}, nil
		}
	}

	return nil, fmt.Errorf("no daytime forecast found")
}

func (s *Service) getForecastURL(pointURL string, headers http.Header) (string, error) {
	req, err := http.NewRequest("GET", pointURL, nil)
	if err != nil {
		return "", err
	}
	req.Header = headers

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data struct {
		Properties struct {
			Forecast string `json:"forecast"`
		} `json:"properties"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.Properties.Forecast, nil
}

func getTemperatureCondition(temp int, unit string) string {
	// Convert to Fahrenheit for consistent comparison if needed
	if unit == "C" {
		temp = int(float64(temp)*9/5 + 32)
	}

	switch {
	case temp >= 85:
		return "hot"
	case temp <= 45:
		return "cold"
	default:
		return "moderate"
	}
}

func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}
