package main

import (
	"github.com/pisarevaa/metrics/internal/server"
	"net/http"
)

func main() {
	storage := server.MemStorage{Metrics: server.MetricGroup{Gauge: make(map[string]float64), Counter: make(map[string]int64)}}
	mux := http.NewServeMux()
	mux.HandleFunc("/", storage.HandleMetrics)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
