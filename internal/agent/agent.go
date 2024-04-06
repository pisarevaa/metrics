package agent

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"runtime"
	"sync"
	"time"
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
		case reflect.Uint64:
			floatValue = float64(value.Uint())
		case reflect.Uint32:
			floatValue = float64(value.Uint())
		case reflect.Float64:
			floatValue = value.Float()
		default:
			return fmt.Errorf("not supported type: %v", value.Kind())
		}
		s.Storage.gauge[v] = floatValue
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
	s.Storage.gauge["RandomValue"] = randomFloat
	return nil
}

func (s *Service) updatePollCount() {
	s.Storage.counter["PollCount"]++
}

func (s *Service) RunUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(s.Config.PollInterval) * time.Second)
	stop := make(chan bool, 1)
	for {
		select {
		case <-ticker.C:
			err := s.UpdateMetrics()
			if err != nil {
				log.Println("error to update metrics:", err)
				stop <- true
			}
		case <-stop:
			return
		}
	}
}

func (s *Service) UpdateMetrics() error {
	log.Println("UpdateMetrics")
	updateMemStatsError := s.updateMemStats()
	if updateMemStatsError != nil {
		return updateMemStatsError
	}
	updateRandomValueError := s.updateRandomValue()
	if updateRandomValueError != nil {
		return updateRandomValueError
	}
	s.updatePollCount()
	return nil
}

func (s *Service) RunSendMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Config.ReportInterval) * time.Second) {
		s.SendMetrics()
	}
}

func (s *Service) SendMetrics() {
	for metric, value := range s.Storage.gauge {
		requestURL := fmt.Sprintf("http://%v/update/gauge/%v/%v", s.Config.Host, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			log.Printf("error making http request: %s\n", err)
		}
	}
	for metric, value := range s.Storage.counter {
		requestURL := fmt.Sprintf("http://%v/update/counter/%v/%v", s.Config.Host, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			log.Printf("error making http request: %s\n", err)
		}
	}
	log.Println("Send Gauge", s.Storage.gauge)
	log.Println("Send Counter", s.Storage.counter)
}
