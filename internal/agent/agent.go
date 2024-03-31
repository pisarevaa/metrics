package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"io"
	"time"
)

const pollIntervalSec = 2
const reportInterval = 10

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

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func sendLog(url string) ([]byte, error) {
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (ms *MemStorage) updateMemStats() {
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
		ms.Gauge[v] = floatValue
	}
}

func (ms *MemStorage) updateRandomValue() {
	ms.Gauge["RandomValue"] = rand.Float64()
}

func (ms *MemStorage) updatePollCount() {
	ms.Counter["PollCount"] += 1
}

func (ms *MemStorage) RunUpdateMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(pollIntervalSec * time.Second) {
		ms.UpdateMetrics()
	}
}

func (ms *MemStorage) UpdateMetrics() {
	ms.updateMemStats()
	ms.updateRandomValue()
	ms.updatePollCount()
}


func (ms MemStorage) RunSendMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	for range time.Tick(reportInterval * time.Second) {
		ms.SendMetrics()
	}
}

func (ms MemStorage) SendMetrics() {
	for metric, value := range ms.Gauge {
		requestURL := fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", metric, value)
		sendLog(requestURL)
	}
	for metric, value := range ms.Counter {
		requestURL := fmt.Sprintf("http://localhost:8080/update/counter/%v/%v", metric, value)
		sendLog(requestURL)
	}
}
