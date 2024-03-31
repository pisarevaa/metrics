package main

import (
	"github.com/pisarevaa/metrics/internal/agent"
	"sync"
)

func main() {
	storage := agent.MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	var wg sync.WaitGroup
	wg.Add(2)
	go storage.RunUpdateMetrics(&wg)
	go storage.RunSendMetrics(&wg)
	wg.Wait()
}
