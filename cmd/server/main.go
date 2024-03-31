package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	metrics MetricGroup
}

type MetricGroup struct {
	gauge  map[string]float64
	counter  map[string]int64
}


func (ms *MemStorage) handleMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	pathParams := strings.Split(strings.TrimPrefix(req.URL.Path, "/update/"), "/")
	if len(pathParams) != 3 {
		http.Error(res, "Path should contains all three fields!", http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := pathParams[0], pathParams[1], pathParams[2]

	if !(metricType == "gauge" || metricType == "counter") {
		http.Error(res, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	if metricName == "" {
		http.Error(res, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	if metricValue == "" || metricValue == "none" {
		http.Error(res, "Empty metricValue is not allowed!", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "metricValue is not corect float!", http.StatusBadRequest)
			return
		}
		ms.metrics.gauge[metricName] = floatValue
	}

	if metricType == "counter" {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "metricValue is not correct integer!", http.StatusBadRequest)
			return
		}
		ms.metrics.counter[metricName] += intValue
	}
}

func main() {
	gauge := make(map[string]float64)
	counter := make(map[string]int64)
	storage := MemStorage{metrics: MetricGroup{gauge, counter}}
	mux := http.NewServeMux()
	mux.HandleFunc("/", storage.handleMetrics)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}