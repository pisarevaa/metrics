package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/pisarevaa/metrics/internal/agent"
	"github.com/pisarevaa/metrics/internal/agent/utils"
)

const processes = 3

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	config := agent.GetConfig()
	client := agent.NewClient()
	logger := agent.GetLogger()
	storage := agent.NewMemStorageRepo()

	semaphore := utils.NewSemaphore(config.RateLimit)
	service := agent.NewService(client, storage, config, logger, semaphore)
	var wg sync.WaitGroup
	wg.Add(processes)
	logger.Info("Client is running...")
	go service.RunUpdateRuntimeMetrics(ctx, &wg)
	go service.RunUpdateGopsutilMetrics(ctx, &wg)
	go service.RunSendMetrics(ctx, &wg)
	wg.Wait()
	logger.Error("exit programm")
}
