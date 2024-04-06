package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/pisarevaa/metrics/internal/storage"
)

type Handler struct {
	Storage *storage.MemStorage
	Config  Config
}

func NewHandler(storage *storage.MemStorage, config Config) *Handler {
	return &Handler{
		Storage: storage,
		Config:  config,
	}
}

func (s *Handler) StoreMetrics(rw http.ResponseWriter, r *http.Request) {
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

	log.Println("Got request", r.URL.Path)
	log.Println("Storage", s.Storage.GetAll())
}

func (s *Handler) GetMetric(rw http.ResponseWriter, r *http.Request) {
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
	log.Println(metricType, metricName)
	value, err := s.Storage.Get(metricType, metricName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	_, errWtrite := io.WriteString(rw, value)
	if errWtrite != nil {
		panic(errWtrite)
	}
}

func (s *Handler) GetAllMetrics(rw http.ResponseWriter, _ *http.Request) {
	metrics := s.Storage.GetAll()
	for key, value := range metrics {
		_, err := io.WriteString(rw, fmt.Sprintf("%v: %v\n", key, value))
		if err != nil {
			panic(err)
		}
	}
}
