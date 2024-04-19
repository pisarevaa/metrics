package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Storage *MemStorage
	Config  Config
	Logger  *zap.SugaredLogger
}

const gauge = "gauge"
const counter = "counter"

func NewHandler(storage *MemStorage, config Config, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		Storage: storage,
		Config:  config,
		Logger:  logger,
	}
}

func (s *Handler) HTTPLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		s.Logger.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}

func (s *Handler) StoreMetrics(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if !(metricType == gauge || metricType == counter) {
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

	if !(metricType == gauge || metricType == counter) {
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
		http.Error(rw, errWtrite.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Handler) GetAllMetrics(rw http.ResponseWriter, _ *http.Request) {
	metrics := s.Storage.GetAll()
	for key, value := range metrics {
		_, err := io.WriteString(rw, fmt.Sprintf("%v: %v\n", key, value))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
