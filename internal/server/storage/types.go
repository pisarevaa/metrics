package storage

import (
	"context"

	"github.com/pisarevaa/metrics/internal/server"
)

type Storage interface {
	GetMetric(ctx context.Context, name string) (metric server.Metrics, err error)
	GetAllMetrics(ctx context.Context) (metric []server.Metrics, err error)
	StoreMetrics(ctx context.Context, metrics []server.Metrics) (err error)
	StoreMetric(ctx context.Context, metric server.Metrics) (err error)
}
