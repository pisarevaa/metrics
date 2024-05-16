package main

import (
	"sync"

	"github.com/pisarevaa/metrics/internal/agent"
)

const processes = 2

func main() {
	config := agent.GetConfig()
	client := agent.NewClient()
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
