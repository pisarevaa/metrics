// Модуль server отвечает за прием метрик от агентов, их хранение и выдачу по запросу.
package main

import (
	"net/http"
	"time"

	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/server/storage"

	_ "net/http/pprof" //nolint:gosec // profiling agent
)

var buildVersion string //nolint:gochecknoglobals // new for task
var buildDate string    //nolint:gochecknoglobals // new for task
var buildCommit string  //nolint:gochecknoglobals // new for task

const readTimeout = 5
const writeTimout = 10
const na = "N/A"

func main() {
	config := server.GetConfig()
	logger := server.GetLogger()

	if buildVersion == "" {
		buildVersion = na
	}
	if buildDate == "" {
		buildDate = na
	}
	if buildCommit == "" {
		buildCommit = na
	}
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
		Handler:      server.MetricsRouter(config, logger, repo),
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimout * time.Second,
	}
	logger.Fatal(srv.ListenAndServe())
}
