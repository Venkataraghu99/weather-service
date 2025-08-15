package api

import (
	"encoding/json"
	"net/http"

	"github.com/Venkataraghu99/weather-service/internal/weather"
)

type Handler struct {
	weatherSvc *weather.Service
}

func NewHandler(weatherSvc *weather.Service) *Handler {
	return &Handler{
		weatherSvc: weatherSvc,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) GetWeatherForecast(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query()
	lat := query.Get("lat")
	long := query.Get("long")

	// Validate parameters
	if lat == "" || long == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "lat and long parameters are required"})
		return
	}

	// Get forecast from weather service
	forecast, err := h.weatherSvc.GetForecast(lat, long)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	// Return successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(forecast)
}
