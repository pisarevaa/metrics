package main

import (
	_ "fmt"
	"net/http"
	"strings"
)

type MemStorage struct {
	metrics []Metric
}

type Metric struct {
	metricName  string
	metricType  string
	metricValue string
}

// Handle metrics logs
func (st *MemStorage) handleMetrics(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	// Get path params
	path_params := strings.Split(strings.TrimPrefix(req.URL.Path, "/update/"), "/")
	if len(path_params) != 3 {
		http.Error(res, "Path should contains all three fields!", http.StatusBadRequest)
		return
	}

	metric := Metric{
		metricName:  path_params[1],
		metricType:  path_params[0],
		metricValue: path_params[2],
	}

	if !(metric.metricType == "gauge" || metric.metricType == "counter") {
		http.Error(res, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	if metric.metricName == "" {
		http.Error(res, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	if metric.metricValue == "" {
		http.Error(res, "Empty metricValue is not allowed!", http.StatusBadRequest)
		return
	}

	st.AddItem(metric)

	res.WriteHeader(http.StatusOK)
}


func (st *MemStorage) AddItem(metric Metric) []Metric {
    st.metrics = append(st.metrics, metric)
    return st.metrics
}

func main() {
	metrics := make([]Metric, 0)
	storage := MemStorage{metrics}
	mux := http.NewServeMux()
	mux.HandleFunc("/", storage.handleMetrics)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}