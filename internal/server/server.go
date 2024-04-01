package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/storage"
	"io"
	"net/http"
)

type Server struct {
	Storage *storage.MemStorage
}

func (s *Server) StoreMetrics(rw http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if !(metricType == "gauge" || metricType == "counter") {
		http.Error(rw, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	if metricName == "" {
		http.Error(rw, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	if metricValue == "" || metricValue == "none" {
		http.Error(rw, "Empty metricValue is not allowed!", http.StatusBadRequest)
		return
	}

	err := s.Storage.Store(metricType, metricName, metricValue)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Got request ", r.URL.Path)
	fmt.Println("Storage Gauge ", s.Storage.Gauge)
	fmt.Println("Storage Counter ", s.Storage.Counter)
}

func (s *Server) GetMetric(rw http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if !(metricType == "gauge" || metricType == "counter") {
		http.Error(rw, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	if metricName == "" {
		http.Error(rw, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	fmt.Println(metricType, metricName)
	value, err := s.Storage.Get(metricType, metricName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	io.WriteString(rw, value)
}

func (s *Server) GetAllMetrics(rw http.ResponseWriter, r *http.Request) {
	metrics := s.Storage.GetAll()
	for key, value := range metrics {
		io.WriteString(rw, fmt.Sprintf("%v: %v\n", key, value))
	}
}
