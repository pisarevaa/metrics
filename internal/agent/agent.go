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

func randomInt() int64 {
	const maxInt = 1000000
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxInt))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	return n
}

func (s *Service) updateMemStats() {
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
			panic("not supported type")
		}
		s.Storage.Gauge[v] = floatValue
	}
}

func (s *Service) updateRandomValue() {
	n1 := randomInt()
	n2 := randomInt()
	randomFloat := float64(n1 / n2)
	s.Storage.Gauge["RandomValue"] = randomFloat
}

func (s *Service) updatePollCount() {
	s.Storage.Counter["PollCount"]++
}

func (s *Service) RunUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Config.PollInterval) * time.Second) {
		s.UpdateMetrics()
	}
}

func (s *Service) UpdateMetrics() {
	log.Println("UpdateMetrics")
	s.updateMemStats()
	s.updateRandomValue()
	s.updatePollCount()
}

func (s *Service) RunSendMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Config.ReportInterval) * time.Second) {
		s.SendMetrics()
	}
}

func (s *Service) SendMetrics() {
	for metric, value := range s.Storage.Gauge {
		requestURL := fmt.Sprintf("http://%v/update/gauge/%v/%v", s.Config.Host, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			log.Printf("error making http request: %s\n", err)
		}
	}
	for metric, value := range s.Storage.Counter {
		requestURL := fmt.Sprintf("http://%v/update/counter/%v/%v", s.Config.Host, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			log.Printf("error making http request: %s\n", err)
		}
	}
	log.Println("Send Gauge ", s.Storage.Gauge)
	log.Println("Send Counter ", s.Storage.Counter)
}
