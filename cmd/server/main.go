package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/storage"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

const readTimeout = 5
const writeTimout = 10

func MetricsRouter(config server.Config) chi.Router {
	storage := storage.NewMemStorageRepo()
	srv := server.NewServer(storage, config)
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Get("/", srv.GetAllMetrics)
	return r
}

func main() {
	config := server.GetConfig()
	log.Printf("Server is running on %v", config.Host)
	srv := &http.Server{
		Addr:         config.Host,
		Handler:      MetricsRouter(config),
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimout * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
