package storage

import (
	"errors"
	"strconv"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (ms *MemStorage) Init() {
	ms.Gauge = make(map[string]float64)
	ms.Counter = make(map[string]int64)
}

func (ms *MemStorage) Store(metricType, metricName, metricValue string) error {
	if metricType == "gauge" {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("metricValue is not corect float")
		}
		ms.Gauge[metricName] = floatValue
	}
	if metricType == "counter" {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("metricValue is not correct integer")
		}
		ms.Counter[metricName] += intValue
	}
	return nil
}

func (ms *MemStorage) Get(metricType, metricName string) (string, error) {
	if metricType == "gauge" {
		value, ok := ms.Gauge[metricName]
		if !ok {
			return "", errors.New("metric is not found")
		}
		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}

	if metricType == "counter" {
		value, ok := ms.Counter[metricName]
		if !ok {
			return "", errors.New("metric is not found")
		}
		return strconv.FormatInt(value, 10), nil
	}

	return "", errors.New("not handled metricType")
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
