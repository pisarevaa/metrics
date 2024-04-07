package main

import (
	"log"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/pisarevaa/metrics/internal/agent"
)

const processes = 2

func main() {
	config := agent.GetConfig()
	client := resty.New()
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(client, storage, config)
	var wg sync.WaitGroup
	wg.Add(processes)
	log.Println("Client is running...")
	go service.RunUpdateMetrics(&wg)
	go service.RunSendMetrics(&wg)
	wg.Wait()
}
