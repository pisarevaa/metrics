package server

import (
	"errors"
	"strconv"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorageRepo() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) Store(metric Metrics) (float64, int64) {
	if metric.MType == gauge {
		if metric.Value == nil {
			ms.gauge[metric.ID] = 0.0
		} else {
			ms.gauge[metric.ID] = *metric.Value
		}
	}
	if metric.MType == counter {
		if metric.Delta != nil {
			ms.counter[metric.ID] += *metric.Delta
		}
	}
	return ms.gauge[metric.ID], ms.counter[metric.ID]
}

func (ms *MemStorage) Get(query QueryMetrics) (*float64, *int64, error) {
	if query.MType == gauge {
		value, ok := ms.gauge[query.ID]
		if !ok {
			return nil, nil, errors.New("metric is not found")
		}
		return &value, nil, nil
	}
	if query.MType == counter {
		value, ok := ms.counter[query.ID]
		if !ok {
			return nil, nil, errors.New("metric is not found")
		}
		return nil, &value, nil
	}
	return nil, nil, errors.New("not handled metricType")
}

func (ms *MemStorage) GetAll() map[string]string {
	metrics := make(map[string]string)
	for key, value := range ms.gauge {
		metrics[key] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	for key, value := range ms.counter {
		metrics[key] = strconv.FormatInt(value, 10)
	}
	return metrics
}
