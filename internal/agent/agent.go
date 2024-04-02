package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

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

type Service struct {
	Client   *resty.Client
	Storage  *MemStorage
	Settings Settings
}

func (s *Service) updateMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
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
	s.Storage.Gauge["RandomValue"] = rand.Float64()
}

func (s *Service) updatePollCount() {
	s.Storage.Counter["PollCount"] += 1
}

func (s *Service) RunUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Settings.pollIntervalSec) * time.Second) {
		s.UpdateMetrics()
	}
}

func (s *Service) UpdateMetrics() {
	fmt.Println("UpdateMetrics")
	s.updateMemStats()
	s.updateRandomValue()
	s.updatePollCount()
}

func (s *Service) RunSendMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(time.Duration(s.Settings.reportInterval) * time.Second) {
		s.SendMetrics()
	}
}

func (s *Service) SendMetrics() {
	for metric, value := range s.Storage.Gauge {
		requestURL := fmt.Sprintf("http://%v/update/gauge/%v/%v", s.Settings.serverHost, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
		}
	}
	for metric, value := range s.Storage.Counter {
		requestURL := fmt.Sprintf("http://%v/update/counter/%v/%v", s.Settings.serverHost, metric, value)
		_, err := s.Client.R().Post(requestURL)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
		}
	}
	fmt.Println("Send Gauge ", s.Storage.Gauge)
	fmt.Println("Send Counter ", s.Storage.Counter)
}
