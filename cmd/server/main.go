package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/storage"
	"net/http"
)

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
	err := http.ListenAndServe(":8080", MetricsRouter())
	if err != nil {
		panic(err)
	}
}
