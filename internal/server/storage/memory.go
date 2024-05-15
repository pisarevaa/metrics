package storage

import (
	"context"
	"errors"
	"sync"
)

type MemStorage struct {
	mx      sync.Mutex
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreMetric(_ context.Context, metric Metrics) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	if metric.MType == Gauge {
		if metric.Value == nil {
			ms.Gauge[metric.ID] = 0.0
		} else {
			ms.Gauge[metric.ID] = *metric.Value
		}
	}
	if metric.MType == Counter {
		if metric.Delta != nil {
			ms.Counter[metric.ID] += *metric.Delta
		}
	}
	return nil
}

func (ms *MemStorage) StoreMetrics(_ context.Context, metrics []Metrics) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	for _, metric := range metrics {
		if metric.MType == Gauge {
			if metric.Value == nil {
				ms.Gauge[metric.ID] = 0.0
			} else {
				ms.Gauge[metric.ID] = *metric.Value
			}
		}
		if metric.MType == Counter {
			if metric.Delta != nil {
				ms.Counter[metric.ID] += *metric.Delta
			}
		}
	}
	return nil
}

func (ms *MemStorage) GetMetric(_ context.Context, name string) (Metrics, error) {
	for metric, value := range ms.Gauge {
		if metric == name {
			return Metrics{
				ID:    metric,
				MType: Gauge,
				Value: &value, // #nosec G601 - проблема ичезнет в go 1.22
			}, nil
		}
	}
	for metric, value := range ms.Counter {
		if metric == name {
			return Metrics{
				ID:    metric,
				MType: Counter,
				Delta: &value, // #nosec G601 - проблема ичезнет в go 1.22
			}, nil
		}
	}
	return Metrics{}, errors.New("metric is not found")
}

func (ms *MemStorage) GetAllMetrics(_ context.Context) ([]Metrics, error) {
	var metrics []Metrics
	for metric, value := range ms.Gauge {
		payload := Metrics{
			ID:    metric,
			MType: Gauge,
			Value: &value, // #nosec G601 - проблема ичезнет в go 1.22
		}
		metrics = append(metrics, payload)
	}
	for metric, value := range ms.Counter {
		payload := Metrics{
			ID:    metric,
			MType: Counter,
			Delta: &value, // #nosec G601 - проблема ичезнет в go 1.22
		}
		metrics = append(metrics, payload)
	}
	return metrics, nil
}

func (ms *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (ms *MemStorage) CloseConnection() {}
