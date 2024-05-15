package agent

import (
	"errors"
	"strconv"
	"sync"
)

type MemStorage struct {
	mx      sync.Mutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorageRepo() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (ms *MemStorage) StoreGauge(metrics map[string]float64) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	for key, value := range metrics {
		ms.Gauge[key] = value
	}
}

func (ms *MemStorage) StoreCounter() {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.Counter["PollCount"]++
}

func (ms *MemStorage) GetMetrics() []Metrics {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	var metrics []Metrics
	for metric, value := range ms.Gauge {
		payload := Metrics{
			ID:    metric,
			MType: gauge,
			Value: value,
		}
		metrics = append(metrics, payload)
	}
	for metric, value := range ms.Counter {
		payload := Metrics{
			ID:    metric,
			MType: counter,
			Delta: value,
		}
		metrics = append(metrics, payload)
	}
	return metrics
}

func (ms *MemStorage) Get(metricType, metricName string) (string, error) {
	if metricType == gauge {
		value, ok := ms.Gauge[metricName]
		if !ok {
			return "", errors.New("metric is not found")
		}
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}

	if metricType == counter {
		value, ok := ms.Counter[metricName]
		if !ok {
			return "", errors.New("metric is not found")
		}
		return strconv.FormatInt(value, 10), nil
	}

	return "", errors.New("not handled metricType")
}

func (ms *MemStorage) GetAll() map[string]string {
	metricsMap := make(map[string]string)
	metrics := ms.GetMetrics()
	for _, metric := range metrics {
		if metric.MType == gauge {
			metricsMap[metric.ID] = strconv.FormatFloat(metric.Value, 'f', -1, 64)
		} else {
			metricsMap[metric.ID] = strconv.FormatInt(metric.Delta, 10)
		}
	}
	return metricsMap
}
