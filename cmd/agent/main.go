package main

import (
	"github.com/pisarevaa/metrics/internal/agent"
	"sync"
)

func main() {
	storage := agent.MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	var wg sync.WaitGroup
	wg.Add(2)
	go storage.UpdateMetrics(&wg)
	go storage.SendMetrics(&wg)
	wg.Wait()
}
