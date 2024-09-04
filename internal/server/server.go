package server

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/pisarevaa/metrics/internal/server/storage"
)

// Создание роутера.
func MetricsRouter(
	ctx context.Context,
	config Config,
	logger *zap.SugaredLogger,
	repo storage.Storage,
) chi.Router {
	if config.Restore {
		metrics, err := LoadFromDosk(config.FileStoragePath)
		if err != nil {
			logger.Error(err)
		}
		err = repo.StoreMetrics(ctx, metrics)
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
	if config.TrustedSubnet != "" {
		r.Use(srv.IPCheckMiddleware)
	}

	r.Mount("/debug", middleware.Profiler())

	r.Get("/ping", srv.Ping)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Post("/update/", srv.StoreMetricsJSON)
	r.Post("/updates/", srv.StoreMetricsJSONBatches)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Post("/value/", srv.GetMetricJSON)
	r.Get("/", srv.GetAllMetrics)

	if config.StoreInterval > 0 {
		logger.Info("Running background tasks...")
		go srv.RunTaskSaveToDisk(ctx)
	}

	return r
}
