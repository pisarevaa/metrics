package main

import (
	"net/http"
	"time"

	"github.com/pisarevaa/metrics/internal/server"
)

const readTimeout = 5
const writeTimout = 10

func main() {
	config := server.GetConfig()
	logger := server.GetLogger()
	storage := server.NewMemStorageRepo()
	logger.Info("Server is running on ", config.Host)
	srv := &http.Server{
		Addr:         config.Host,
		Handler:      server.MetricsRouter(config, logger, storage),
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimout * time.Second,
	}
	logger.Fatal(srv.ListenAndServe())
}
