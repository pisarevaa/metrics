package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pisarevaa/metrics/internal/agent"
)

const processes = 3

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := agent.GetConfig()
	client := agent.NewClient()
	logger := agent.GetLogger()
	storage := agent.NewMemStorageRepo()

	go func() {
		s := <-sig
		logger.Info("Received signal:", s.String())
		cancel()
	}()

	semaphore := agent.NewSemaphore(config.RateLimit)
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
