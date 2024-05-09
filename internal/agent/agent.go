package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	ID    string  `json:"id"`    // имя метрики
	MType string  `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta"` // значение метрики в случае передачи counter
	Value float64 `json:"value"` // значение метрики в случае передачи gauge
}

const (
	gauge   = "gauge"
	counter = "counter"
)

func randomInt() (int64, error) {
	const maxInt = 1000000
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxInt))
	if err != nil {
		return 0, err
	}
	n := nBig.Int64()
	return n, nil
}

func (s *Service) updateMemStats() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var gaugeMetrics = [...]string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}

	for _, v := range gaugeMetrics {
		value := reflect.ValueOf(memStats).FieldByName(v)
		var floatValue float64
		switch value.Kind() {
		case reflect.Uint64, reflect.Uint32:
			floatValue = float64(value.Uint())
		case reflect.Float64:
			floatValue = value.Float()
		default:
			return fmt.Errorf("not supported type: %v", value.Kind())
		}
		s.Storage.StoreGauge(v, floatValue)
	}
	return nil
}

func (s *Service) updateRandomValue() error {
	n1, err1 := randomInt()
	if err1 != nil {
		return err1
	}
	n2, err2 := randomInt()
	if err2 != nil {
		return err2
	}
	randomFloat := float64(n1 / n2)
	s.Storage.StoreGauge("RandomValue", randomFloat)
	return nil
}

func (s *Service) RunUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(s.Config.PollInterval) * time.Second)
	defer ticker.Stop()
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			err := s.UpdateMetrics()
			if err != nil {
				s.Logger.Error("error to update metrics:", err)
				stop <- true
			}
		case <-stop:
			return
		}
	}
}

func (s *Service) UpdateMetrics() error {
	s.Logger.Info("UpdateMetrics")
	updateMemStatsError := s.updateMemStats()
	if updateMemStatsError != nil {
		return updateMemStatsError
	}
	updateRandomValueError := s.updateRandomValue()
	if updateRandomValueError != nil {
		return updateRandomValueError
	}
	s.Storage.StoreCounter()
	return nil
}

func (s *Service) RunSendMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Config.ReportInterval) * time.Second) {
		s.SendMetrics()
	}
}

// func (s *Service) makeHTTPRequest(metric Metrics) {
// 	requestURL := fmt.Sprintf("http://%v/update/", s.Config.Host)
// 	buf := bytes.NewBuffer(nil)
// 	zb := gzip.NewWriter(buf)
// 	payloadString, err := json.Marshal(metric)
// 	if err != nil {
// 		s.Logger.Error(err)
// 		return
// 	}
// 	_, err = zb.Write(payloadString)
// 	if err != nil {
// 		s.Logger.Error(err)
// 		return
// 	}
// 	err = zb.Close()
// 	if err != nil {
// 		s.Logger.Error(err)
// 		return
// 	}
// 	_, err = s.Client.R().
// 		SetHeader("Content-Type", "application/json").
// 		SetHeader("Content-Encoding", "gzip").
// 		SetBody(buf).
// 		Post(requestURL)
// 	if err != nil {
// 		s.Logger.Error("error making http request: ", err)
// 	}
// }

func (s *Service) makeHTTPRequest(metrics []Metrics) {
	requestURL := fmt.Sprintf("http://%v/updates/", s.Config.Host)
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	payloadString, err := json.Marshal(metrics)
	if err != nil {
		s.Logger.Error(err)
		return
	}
	_, err = zb.Write(payloadString)
	if err != nil {
		s.Logger.Error(err)
		return
	}
	err = zb.Close()
	if err != nil {
		s.Logger.Error(err)
		return
	}
	_, err = s.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post(requestURL)
	if err != nil {
		s.Logger.Error("error making http request: ", err)
	}
}

// func (s *Service) SendMetrics() {
// 	metrics := s.Storage.GetMetrics()
// 	for _, metric := range metrics {
// 		s.makeHTTPRequest(metric)
// 	}
// 	s.Logger.Info("Send Gauge", s.Storage.Gauge)
// 	s.Logger.Info("Send Counter", s.Storage.Counter)
// }

func (s *Service) SendMetrics() {
	metrics := s.Storage.GetMetrics()
	s.makeHTTPRequest(metrics)
	s.Logger.Info("Send Gauge", s.Storage.Gauge)
	s.Logger.Info("Send Counter", s.Storage.Counter)
}
