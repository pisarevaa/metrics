package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func MetricsRouter(config Config, logger *zap.SugaredLogger) chi.Router {
	storage := NewMemStorageRepo()
	srv := NewHandler(storage, config, logger)
	r := chi.NewRouter()
	r.Use(srv.HTTPLoggingMiddleware)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Post("/update/", srv.StoreMetricsJSON)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Post("/value/", srv.GetMetricJSON)
	r.Get("/", srv.GetAllMetrics)
	return r
}
