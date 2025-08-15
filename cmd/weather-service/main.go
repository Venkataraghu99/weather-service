package main

import (
	"log"
	"net/http"

	"github.com/Venkataraghu99/weather-service/internal/api"
	"github.com/Venkataraghu99/weather-service/internal/weather"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the weather service
	weatherSvc := weather.NewService()

	// Initialize the API handler
	handler := api.NewHandler(weatherSvc)

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/api/weather/forecast", handler.GetWeatherForecast).Methods("GET")

	// Start the server
	port := ":8080"
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
