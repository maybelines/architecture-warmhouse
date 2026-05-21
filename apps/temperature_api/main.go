package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func randomTemp() float64 {
	return math.Round((15.0+rand.Float64()*20.0)*10) / 10
}

func handleTemperature(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/temperature")
	path = strings.TrimPrefix(path, "/")

	var sensorID, location string

	if path == "" {
		location = r.URL.Query().Get("location")
	} else {
		sensorID = path
	}

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	resp := TemperatureResponse{
		Value:       randomTemp(),
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Temperature at %s", location),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + strings.TrimPrefix(port, ":")

	mux := http.NewServeMux()
	mux.HandleFunc("/temperature", handleTemperature)
	mux.HandleFunc("/temperature/", handleTemperature)
	mux.HandleFunc("/health", handleHealth)

	log.Printf("Temperature API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
