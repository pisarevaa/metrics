package storage

import (
	"context"
)

type Storage interface {
	GetMetric(ctx context.Context, id string, mtype string) (metric Metrics, err error)
	GetAllMetrics(ctx context.Context) (metric []Metrics, err error)
	StoreMetrics(ctx context.Context, metrics []Metrics) (err error)
	StoreMetric(ctx context.Context, metric Metrics) (err error)
	Ping(ctx context.Context) (err error)
	CloseConnection()
}

type Metrics struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type QueryMetrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)
