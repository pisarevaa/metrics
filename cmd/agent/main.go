// Модуль agent отвечает за отправку метрик на сервер.
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"net/http"

	"github.com/pisarevaa/metrics/internal/agent"
	"github.com/pisarevaa/metrics/internal/agent/utils"

	_ "net/http/pprof" //nolint:gosec // profiling agent
)

var buildVersion, buildDate, buildCommit string //nolint:gochecknoglobals // new for task

const processes = 3 // количество гоурутин
const readTimeout = 5
const writeTimout = 10

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	config := agent.GetConfig()
	client := agent.NewClient()
	logger := agent.GetLogger()
	storage := agent.NewMemStorageRepo()

	utils.SetDefaultBuildInfo(&buildVersion)
	utils.SetDefaultBuildInfo(&buildDate)
	utils.SetDefaultBuildInfo(&buildCommit)
	logger.Info("Build version: ", buildVersion)
	logger.Info("Build date: ", buildDate)
	logger.Info("Build commit: ", buildCommit)

	semaphore := utils.NewSemaphore(config.RateLimit)
	service := agent.NewService(client, storage, config, logger, semaphore)

	// Profiling agent http://127.0.0.1:8080/debug/pprof/
	httpServer := &http.Server{
		Addr:         "localhost:8085",
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimout * time.Second,
	}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			_ = httpServer.Shutdown(context.Background())
		}
	}()

	var wg sync.WaitGroup
	wg.Add(processes)
	logger.Info("Client is running...")
	go service.RunUpdateRuntimeMetrics(ctx, &wg)
	go service.RunUpdateGopsutilMetrics(ctx, &wg)
	go service.RunSendMetrics(ctx, &wg)
	wg.Wait()

	logger.Error("exit programm")
}
