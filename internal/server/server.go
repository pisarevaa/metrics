package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server/storage"
	"go.uber.org/zap"
)

func MetricsRouter(
	config Config,
	logger *zap.SugaredLogger,
	repo storage.Storage,
) chi.Router {
	if config.Restore {
		metrics, err := LoadFromDosk(config.FileStoragePath)
		if err != nil {
			logger.Error(err)
		}
		err = repo.StoreMetrics(context.Background(), metrics)
		if err != nil {
			logger.Error(err)
		}
	}
	srv := NewHandler(config, logger, repo)
	r := chi.NewRouter()
	r.Use(srv.HTTPLoggingMiddleware)
	r.Use(srv.GzipMiddleware)

	if config.Key != "" {
		r.Use(srv.HashCheckMiddleware)
	}
	r.Get("/ping", srv.Ping)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Post("/update/", srv.StoreMetricsJSON)
	r.Post("/updates/", srv.StoreMetricsJSONBatches)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Post("/value/", srv.GetMetricJSON)
	r.Get("/", srv.GetAllMetrics)

	if config.StoreInterval > 0 {
		logger.Info("Running background tasks...")
		go srv.RunTaskSaveToDisk()
	}

	return r
}
