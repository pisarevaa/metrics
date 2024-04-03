package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/storage"
	"net/http"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

func MetricsRouter(config server.Config) chi.Router {
	storage := storage.MemStorage{}
	storage.Init()
	server := server.Server{Storage: &storage, Config: config}
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", server.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", server.GetMetric)
	r.Get("/", server.GetAllMetrics)
	return r
}

func main() {
	config := server.GetConfigs()
	fmt.Printf("Server is running on %v", config.Host)
	err := http.ListenAndServe(config.Host, MetricsRouter(config))
	if err != nil {
		panic(err)
	}
}
