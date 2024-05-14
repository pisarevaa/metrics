package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/pisarevaa/metrics/internal/server"
)

type MemStorage struct {
	mx      sync.Mutex
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

func NewMemStorageRepo() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreMetric(_ context.Context, metric server.Metrics) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	if metric.MType == server.Gauge {
		if metric.Value == nil {
			ms.Gauge[metric.ID] = 0.0
		} else {
			ms.Gauge[metric.ID] = *metric.Value
		}
	}
	if metric.MType == server.Counter {
		if metric.Delta != nil {
			ms.Counter[metric.ID] += *metric.Delta
		}
	}
	return nil
}

func (ms *MemStorage) StoreMetrics(_ context.Context, metrics []server.Metrics) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	for _, metric := range metrics {
		if metric.MType == server.Gauge {
			if metric.Value == nil {
				ms.Gauge[metric.ID] = 0.0
			} else {
				ms.Gauge[metric.ID] = *metric.Value
			}
		}
		if metric.MType == server.Counter {
			if metric.Delta != nil {
				ms.Counter[metric.ID] += *metric.Delta
			}
		}
	}
	return nil
}

func (ms *MemStorage) GetMetric(_ context.Context, name string) (server.Metrics, error) {
	for metric, value := range ms.Gauge {
		if metric == name {
			return server.Metrics{
				ID:    metric,
				MType: server.Gauge,
				Value: &value, // #nosec G601 - проблема ичезнет в go 1.22
			}, nil
		}
	}
	for metric, value := range ms.Counter {
		if metric == name {
			return server.Metrics{
				ID:    metric,
				MType: server.Counter,
				Delta: &value, // #nosec G601 - проблема ичезнет в go 1.22
			}, nil
		}
	}
	return server.Metrics{}, errors.New("metric is not found")
}

func (ms *MemStorage) GetAllMetrics(_ context.Context) ([]server.Metrics, error) {
	var metrics []server.Metrics
	for metric, value := range ms.Gauge {
		payload := server.Metrics{
			ID:    metric,
			MType: server.Gauge,
			Value: &value, // #nosec G601 - проблема ичезнет в go 1.22
		}
		metrics = append(metrics, payload)
	}
	for metric, value := range ms.Counter {
		payload := server.Metrics{
			ID:    metric,
			MType: server.Counter,
			Delta: &value, // #nosec G601 - проблема ичезнет в go 1.22
		}
		metrics = append(metrics, payload)
	}
	return metrics, nil
}
