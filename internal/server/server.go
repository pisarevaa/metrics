package server

import (
	"github.com/go-chi/chi/v5"
)

func MetricsRouter(config Config) chi.Router {
	storage := NewMemStorageRepo()
	srv := NewHandler(storage, config)
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Get("/", srv.GetAllMetrics)
	return r
}
