package server

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
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

func (ms *MemStorage) Store(metric Metrics) (float64, int64) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	if metric.MType == gauge {
		if metric.Value == nil {
			ms.Gauge[metric.ID] = 0.0
		} else {
			ms.Gauge[metric.ID] = *metric.Value
		}
	}
	if metric.MType == counter {
		if metric.Delta != nil {
			ms.Counter[metric.ID] += *metric.Delta
		}
	}
	return ms.Gauge[metric.ID], ms.Counter[metric.ID]
}

func (ms *MemStorage) SaveToDosk(filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&ms)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ms *MemStorage) LoadFromDosk(filename string) error {
	if filename == "" {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(ms)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ms *MemStorage) Get(query QueryMetrics) (*float64, *int64, error) {
	if query.MType == gauge {
		value, ok := ms.Gauge[query.ID]
		if !ok {
			return nil, nil, errors.New("metric is not found")
		}
		return &value, nil, nil
	}
	if query.MType == counter {
		value, ok := ms.Counter[query.ID]
		if !ok {
			return nil, nil, errors.New("metric is not found")
		}
		return nil, &value, nil
	}
	return nil, nil, errors.New("not handled metricType")
}

func (ms *MemStorage) GetAll() map[string]string {
	metrics := make(map[string]string)
	for key, value := range ms.Gauge {
		metrics[key] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	for key, value := range ms.Counter {
		metrics[key] = strconv.FormatInt(value, 10)
	}
	return metrics
}
