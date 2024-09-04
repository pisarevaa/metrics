// Модуль server отвечает за прием метрик от агентов, их хранение и выдачу по запросу.
package main

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/server/storage"
	"github.com/pisarevaa/metrics/internal/server/utils"

	_ "net/http/pprof" //nolint:gosec // profiling agent

	pb "github.com/pisarevaa/metrics/proto"
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

	var grpcServer *grpc.Server
	if config.GrpcActive {
		listen, err := net.Listen("tcp", ":"+config.GrpcPort)
		if err != nil {
			logger.Fatal(err)
		}
		grpcServer = grpc.NewServer()
		pb.RegisterMetricsServer(grpcServer, server.NewGrpcServer(config, logger, repo))
		logger.Info("gRPC server is running...")
		go func() {
			if errGrpc := grpcServer.Serve(listen); errGrpc != nil {
				logger.Info("Could not listen on tcp:" + config.GrpcPort)
			}
		}()
	}

	go func() {
		if errServer := srv.ListenAndServe(); errServer != nil && errServer != http.ErrServerClosed {
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
	if config.GrpcActive {
		grpcServer.GracefulStop()
	}
	logger.Info("Server is gracefully shutdown")
}
