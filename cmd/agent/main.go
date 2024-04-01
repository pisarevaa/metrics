package main

import (
	"github.com/pisarevaa/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
	"sync"
)

func main() {
	client := resty.New()
	storage := agent.MemStorage{}
	storage.Init()
	service := agent.Service{Storage: &storage, Client: client}
	var wg sync.WaitGroup
	wg.Add(2)
	go service.RunUpdateMetrics(&wg)
	go service.RunSendMetrics(&wg)
	wg.Wait()
}
