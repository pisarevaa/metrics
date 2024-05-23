package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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

func (s *Service) updateMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var gaugeMetrics = map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
	}
	s.Storage.StoreGauge(gaugeMetrics)
}

func (s *Service) updateGopsutilStats() error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	counts, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	var gaugeMetrics = map[string]float64{
		"TotalMemory":     float64(v.Total),
		"FreeMemory":      float64(v.Free),
		"CPUutilization1": counts[0],
	}
	s.Storage.StoreGauge(gaugeMetrics)
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
	s.Storage.StoreGauge(map[string]float64{"RandomValue": float64(n1 / n2)})
	return nil
}

func (s *Service) RunUpdateGopsutilMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(s.Config.PollInterval) * time.Second)
	defer ticker.Stop()
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			err := s.updateGopsutilStats()
			if err != nil {
				s.Logger.Error("error to update metrics:", err)
				stop <- true
			}
		case <-stop:
			return
		}
	}
}

func (s *Service) RunUpdateRuntimeMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(s.Config.PollInterval) * time.Second)
	defer ticker.Stop()
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			err := s.UpdateRuntimeMetrics()
			if err != nil {
				s.Logger.Error("error to update metrics:", err)
				stop <- true
			}
		case <-stop:
			return
		}
	}
}

func (s *Service) UpdateRuntimeMetrics() error {
	s.Logger.Info("UpdateRuntimeMetrics")
	s.updateMemStats()
	updateRandomValueError := s.updateRandomValue()
	if updateRandomValueError != nil {
		return updateRandomValueError
	}
	s.Storage.StoreCounter()
	return nil
}

func (s *Service) UpdateGopsutilMetrics() error {
	s.Logger.Info("UpdateGopsutilMetrics")
	s.updateMemStats()
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
	r := s.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf)
	if s.Config.Key != "" {
		hash, errHash := GetBodyHash(payloadString, s.Config.Key)
		if errHash != nil {
			s.Logger.Error(errHash)
			return
		}
		s.Logger.Info("Hash", hash)
		r.SetHeader("Hash", hash)
	}
	_, err = r.Post(requestURL)
	if err != nil {
		s.Logger.Error("error making http request: ", err)
	}
}

func (s *Service) SendMetrics() {
	metrics := s.Storage.GetMetrics()
	s.makeHTTPRequest(metrics)
	s.Logger.Info("Send Gauge", s.Storage.Gauge)
	s.Logger.Info("Send Counter", s.Storage.Counter)
}
