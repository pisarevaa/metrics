package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server/storage"
	"go.uber.org/zap"
)

type Handler struct {
	Config  Config
	Logger  *zap.SugaredLogger
	Storage storage.Storage
}

type QueryMetrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

func NewHandler(config Config, logger *zap.SugaredLogger, repo storage.Storage) *Handler {
	return &Handler{
		Config:  config,
		Logger:  logger,
		Storage: repo,
	}
}

func (s *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := s.Storage.Ping(r.Context())
	if err != nil {
		http.Error(w, "DBPool is not initialized!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) StoreMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if !(metricType == storage.Gauge || metricType == storage.Counter) {
		http.Error(w, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	if metricName == "" {
		http.Error(w, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	if metricValue == "" || metricValue == "none" {
		http.Error(w, "Empty metricValue is not allowed!", http.StatusBadRequest)
		return
	}

	metric := storage.Metrics{
		ID:    metricName,
		MType: metricType,
	}

	if metricType == storage.Gauge {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "metricValue is not corect float", http.StatusBadRequest)
			return
		}
		metric.Value = floatValue
	}
	if metricType == storage.Counter {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "metricValue is not correct integer", http.StatusBadRequest)
			return
		}
		metric.Delta = intValue
	}

	err := s.Storage.StoreMetric(r.Context(), metric)
	if err != nil {
		http.Error(w, "Error to store metric", http.StatusBadRequest)
		return
	}

	s.Logger.Info("Got request ", r.URL.Path)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) StoreMetricsJSON(w http.ResponseWriter, r *http.Request) { //nolint:funlen /// TODO: refactor
	var metric storage.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Logger.Info("metric ", metric)

	if !(metric.MType == storage.Gauge || metric.MType == storage.Counter) {
		s.Logger.Error("Only 'gauge' and 'counter' values are allowed!")
		http.Error(w, "Only 'gauge' and 'counter' values are allowed!", http.StatusBadRequest)
		return
	}
	if metric.ID == "" {
		s.Logger.Error("Empty metric id is not allowed!")
		http.Error(w, "Empty metric id is not allowed!", http.StatusNotFound)
		return
	}

	err = s.Storage.StoreMetric(r.Context(), metric)
	if err != nil {
		s.Logger.Error("Error to store metric ", err)
		http.Error(w, "Error to store metric", http.StatusBadRequest)
		return
	}

	if s.Config.StoreInterval == 0 {
		metrics, errMetrics := s.Storage.GetAllMetrics(r.Context())
		if errMetrics != nil {
			http.Error(w, "Error to get all metrics", http.StatusBadRequest)
			return
		}
		err = SaveToDisk(metrics, s.Config.FileStoragePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	metricStored, err := s.Storage.GetMetric(r.Context(), metric.ID, metric.MType)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := metricStored.ToJSON()
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Logger.Info("Got request ", r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) StoreMetricsJSONBatches(w http.ResponseWriter, r *http.Request) {
	var metrics []storage.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, metric := range metrics {
		if !(metric.MType == storage.Gauge || metric.MType == storage.Counter) {
			http.Error(w, "Only 'gauge' and 'counter' values are allowed!", http.StatusBadRequest)
			return
		}
		if metric.ID == "" {
			http.Error(w, "Empty metric id is not allowed!", http.StatusNotFound)
			return
		}
	}

	err = s.Storage.StoreMetrics(r.Context(), metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.Config.StoreInterval == 0 {
		err = SaveToDisk(metrics, s.Config.FileStoragePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.Logger.Info("Got request ", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if !(metricType == storage.Gauge || metricType == storage.Counter) {
		http.Error(w, "Only 'gauge' and 'counter' values are allowed!", http.StatusBadRequest)
		return
	}
	if metricName == "" {
		http.Error(w, "Empty metricName is not allowed!", http.StatusNotFound)
		return
	}
	s.Logger.Info(metricType, metricName)

	query := QueryMetrics{
		ID:    metricName,
		MType: metricType,
	}

	metric, err := s.Storage.GetMetric(r.Context(), query.ID, query.MType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if query.MType == storage.Gauge {
		valueString := strconv.FormatFloat(metric.Value, 'f', -1, 64)
		_, errWtrite := io.WriteString(w, valueString)
		if errWtrite != nil {
			http.Error(w, errWtrite.Error(), http.StatusBadRequest)
			return
		}
	}
	if query.MType == storage.Counter {
		valueString := strconv.FormatInt(metric.Delta, 10)
		_, errWtrite := io.WriteString(w, valueString)
		if errWtrite != nil {
			http.Error(w, errWtrite.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var query QueryMetrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &query); err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(query.MType == storage.Gauge || query.MType == storage.Counter) {
		s.Logger.Error("Only 'gauge' and 'counter' values are not allowed!")
		http.Error(w, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	s.Logger.Info(query.MType, query.ID)

	metric, err := s.Storage.GetMetric(r.Context(), query.ID, query.MType)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := metric.ToJSON()
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		s.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := s.Storage.GetAllMetrics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for _, value := range metrics {
		var row string
		if value.MType == storage.Gauge {
			row = fmt.Sprintf("%v: %v\n", value.ID, value.Value)
		} else {
			row = fmt.Sprintf("%v: %v\n", value.ID, value.Delta)
		}
		_, err = w.Write([]byte(row))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Handler) RunTaskSaveToDisk() {
	ticker := time.NewTicker(time.Duration(s.Config.StoreInterval) * time.Second)
	defer ticker.Stop()
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			metrics, err := s.Storage.GetAllMetrics(context.Background())
			if err != nil {
				s.Logger.Error("error to save metrics to disk:", err)
				stop <- true
			}
			err = SaveToDisk(metrics, s.Config.FileStoragePath)
			if err != nil {
				s.Logger.Error("error to save metrics to disk:", err)
				stop <- true
			}
			s.Logger.Info("success save metrics to disk")
		case <-stop:
			return
		}
	}
}
