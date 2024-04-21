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
	r.Post("/update/", srv.StoreMetrics)
	r.Get("/value/", srv.GetMetric)
	r.Get("/", srv.GetAllMetrics)
	return r
}
