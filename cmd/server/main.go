// Модуль server отвечает за прием метрик от агентов, их хранение и выдачу по запросу.
package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/server/storage"
	"github.com/pisarevaa/metrics/internal/server/utils"

	_ "net/http/pprof" //nolint:gosec // profiling agent
)

var buildVersion, buildDate, buildCommit string //nolint:gochecknoglobals // new for task

const readTimeout = 5
const writeTimeout = 10
const shutdownTimeout = 10

func main() {
	ctxCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctxStop, stop := signal.NotifyContext(ctxCancel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	config := server.GetConfig()
	logger := server.GetLogger()

	utils.SetDefaultBuildInfo(&buildVersion)
	utils.SetDefaultBuildInfo(&buildDate)
	utils.SetDefaultBuildInfo(&buildCommit)
	logger.Info("Build version: ", buildVersion)
	logger.Info("Build date: ", buildDate)
	logger.Info("Build commit: ", buildCommit)

	var repo storage.Storage
	if config.DatabaseDSN == "" {
		repo = storage.NewMemStorage()
	} else {
		repo = storage.NewDBStorage(config.DatabaseDSN, logger)
	}
	defer repo.CloseConnection()
	logger.Info("Server is running on ", config.Host)
	srv := &http.Server{
		Addr:         config.Host,
		Handler:      server.MetricsRouter(ctxStop, config, logger, repo),
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Info("Could not listen on ", config.Host)
		}
	}()
	<-ctxStop.Done()
	shutdownCtx, timeout := context.WithTimeout(ctxStop, shutdownTimeout*time.Second)
	defer timeout()
	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Server is gracefully shutdown")
}
