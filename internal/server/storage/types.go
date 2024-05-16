package storage

import (
	"context"
	"encoding/json"
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

type GaugeMetrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Value float64 `json:"value"`
}

type CounterMetrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Delta int64  `json:"delta"`
}

type QueryMetrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func (m *Metrics) ToJSON() ([]byte, error) {
	var resp []byte
	var err error
	if m.MType == Gauge {
		resp, err = json.Marshal(GaugeMetrics{ID: m.ID, MType: m.MType, Value: m.Value})
	}
	if m.MType == Counter {
		resp, err = json.Marshal(CounterMetrics{ID: m.ID, MType: m.MType, Delta: m.Delta})
	}
	return resp, err
}
