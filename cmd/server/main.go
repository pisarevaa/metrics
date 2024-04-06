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
