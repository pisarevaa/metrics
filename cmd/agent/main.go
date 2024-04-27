package main

import (
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/pisarevaa/metrics/internal/agent"
)

const processes = 2

func main() {
	config := agent.GetConfig()
	client := resty.New()
	logger := agent.GetLogger()
	storage := agent.NewMemStorageRepo()
	service := agent.NewService(client, storage, config, logger)
	var wg sync.WaitGroup
	wg.Add(processes)
	logger.Info("Client is running...")
	go service.RunUpdateMetrics(&wg)
	go service.RunSendMetrics(&wg)
	wg.Wait()
}
