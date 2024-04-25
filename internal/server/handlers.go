package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	Storage *MemStorage
	Config  Config
	Logger  *zap.SugaredLogger
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type QueryMetrics struct {
	ID    string `json:"id"`
	MType string `json:"type"`
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

func (s *Handler) StoreMetrics(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if !(metricType == gauge || metricType == counter) {
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

	metric := Metrics{
		ID:    metricName,
		MType: metricType,
	}

	if metricType == gauge {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "metricValue is not corect float", http.StatusBadRequest)
			return
		}
		metric.Value = &floatValue
	}
	if metricType == counter {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "metricValue is not correct integer", http.StatusBadRequest)
			return
		}
		metric.Delta = &intValue
	}

	s.Storage.Store(metric)

	s.Logger.Info("Got request ", r.URL.Path)
	s.Logger.Info("Storage ", s.Storage.GetAll())

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) StoreMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var metric Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(metric.MType == gauge || metric.MType == counter) {
		http.Error(w, "Only 'gauge' and 'counter' values are allowed!", http.StatusBadRequest)
		return
	}
	if metric.ID == "" {
		http.Error(w, "Empty metric id is not allowed!", http.StatusNotFound)
		return
	}
	if metric.MType == gauge && metric.Value == nil {
		http.Error(w, "Empty metric Value is not allowed!", http.StatusBadRequest)
		return
	}
	if metric.MType == counter && metric.Delta == nil {
		http.Error(w, "Empty metric Delta is not allowed!", http.StatusBadRequest)
		return
	}

	value, delta := s.Storage.Store(metric)

	if s.Config.StoreInterval == 0 {
		err = s.Storage.SaveToDosk(s.Config.FileStoragePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(Metrics{
		ID:    metric.ID,
		MType: metric.MType,
		Delta: &delta,
		Value: &value,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Logger.Info("Got request ", r.URL.Path)
	s.Logger.Info("Storage ", s.Storage.GetAll())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if !(metricType == gauge || metricType == counter) {
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

	value, delta, err := s.Storage.Get(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if metricType == gauge && value != nil {
		valueString := strconv.FormatFloat(*value, 'f', -1, 64)
		_, errWtrite := io.WriteString(w, valueString)
		if errWtrite != nil {
			http.Error(w, errWtrite.Error(), http.StatusBadRequest)
			return
		}
	}
	if metricType == counter && delta != nil {
		valueString := strconv.FormatInt(*delta, 10)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !(query.MType == gauge || query.MType == counter) {
		http.Error(w, "Only 'gauge' and 'counter' values are not allowed!", http.StatusBadRequest)
		return
	}
	// if query.ID == "" {
	// 	http.Error(w, "Empty id is not allowed!", http.StatusNotFound)
	// 	return
	// }
	s.Logger.Info(query.MType, query.ID)

	value, delta, err := s.Storage.Get(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	resp, err := json.Marshal(Metrics{
		ID:    query.ID,
		MType: query.MType,
		Delta: delta,
		Value: value,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Handler) GetAllMetrics(w http.ResponseWriter, _ *http.Request) {
	metrics := s.Storage.GetAll()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for key, value := range metrics {
		row := fmt.Sprintf("%v: %v\n", key, value)
		_, err := w.Write([]byte(row))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s *Handler) RunTaskSaveToDisk() {
	ticker := time.NewTicker(time.Duration(s.Config.StoreInterval) * time.Second)
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			err := s.Storage.SaveToDosk(s.Config.FileStoragePath)
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
