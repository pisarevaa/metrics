package server

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func MetricsRouter(config Config, logger *zap.SugaredLogger, storage *MemStorage) chi.Router {
	if config.Restore {
		err := storage.LoadFromDosk(config.FileStoragePath)
		if err != nil {
			logger.Error(err)
		}
	}
	srv := NewHandler(storage, config, logger)
	r := chi.NewRouter()
	r.Use(srv.HTTPLoggingMiddleware)
	r.Use(srv.GzipMiddleware)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", srv.StoreMetrics)
	r.Post("/update/", srv.StoreMetricsJSON)
	r.Get("/value/{metricType}/{metricName}", srv.GetMetric)
	r.Post("/value/", srv.GetMetricJSON)
	r.Get("/", srv.GetAllMetrics)

	if config.StoreInterval > 0 {
		logger.Info("Running background tasks...")
		go srv.RunTaskSaveToDisk()
	}

	return r
}
