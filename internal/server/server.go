package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/pisarevaa/metrics/internal/storage"
)

func MetricsRouter(config Config) chi.Router {
	storage := storage.NewMemStorageRepo()
	srv := NewHandler(storage, config)
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Get("/", srv.GetAllMetrics)
	return r
}
