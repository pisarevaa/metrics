package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/storage"
)

var host string

func MetricsRouter() chi.Router {
	storage := storage.MemStorage{}
	storage.Init()
	server := server.Server{Storage: &storage}
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", server.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", server.GetMetric)
	r.Get("/", server.GetAllMetrics)
	return r
}

func main() {
	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	fmt.Printf("Server is running on %v", host)
	err := http.ListenAndServe(host, MetricsRouter())
	if err != nil {
		panic(err)
	}
}
