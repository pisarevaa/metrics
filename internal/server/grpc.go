package server

import (
	"context"

	"go.uber.org/zap"

	"github.com/pisarevaa/metrics/internal/server/storage"
	pb "github.com/pisarevaa/metrics/proto"
)

type GrpcServer struct {
	pb.UnimplementedMetricsServer
	Config  Config
	Logger  *zap.SugaredLogger
	Storage storage.Storage
}

// Создание сервера GRPC.
func NewGrpcServer(config Config, logger *zap.SugaredLogger, repo storage.Storage) *GrpcServer {
	return &GrpcServer{
		Config:  config,
		Logger:  logger,
		Storage: repo,
	}
}

func (g *GrpcServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse
	var mType string
	switch in.GetType() {
	case pb.MetricType_gauge:
		mType = storage.Gauge
	case pb.MetricType_counter:
		mType = storage.Counter
	}
	metric, err := g.Storage.GetMetric(ctx, in.GetId(), mType)

	if err != nil {
		g.Logger.Error(err)
		response.Error = err.Error()
		return &response, nil
	}

	g.Logger.Info("metric: ", metric)

	response.Metric.Id = in.GetId()
	response.Metric.Type = in.GetType()
	response.Metric.Value = metric.Value
	response.Metric.Delta = metric.Delta

	return &response, nil
}

func (g *GrpcServer) GetMetrics(ctx context.Context, _ *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	var response pb.GetMetricsResponse

	metrics, err := g.Storage.GetAllMetrics(ctx)
	if err != nil {
		g.Logger.Error(err)
		response.Error = err.Error()
		return &response, nil
	}

	for _, value := range metrics {
		var mType pb.MetricType
		switch value.MType {
		case storage.Gauge:
			mType = 0
		case storage.Counter:
			mType = 1
		}
		response.Metrics = append(response.Metrics, &pb.Metric{
			Id:    value.ID,
			Type:  mType,
			Value: value.Value,
			Delta: value.Delta,
		})
	}

	return &response, nil
}

func (g *GrpcServer) AddMetrics(ctx context.Context, in *pb.AddMetricsRequest) (*pb.AddMetricsResponse, error) {
	var metrics []storage.Metrics
	var response pb.AddMetricsResponse

	for _, metric := range in.GetMetrics() {
		var mType string
		switch metric.GetType() {
		case pb.MetricType_gauge:
			mType = storage.Gauge
		case pb.MetricType_counter:
			mType = storage.Counter
		}
		metrics = append(metrics, storage.Metrics{
			ID:    metric.GetId(),
			MType: mType,
			Value: metric.GetValue(),
			Delta: metric.GetDelta(),
		})
	}

	g.Logger.Info("metrics updates: ", metrics)

	err := g.Storage.StoreMetrics(ctx, metrics)
	if err != nil {
		g.Logger.Error(err)
		response.Error = err.Error()
		return &response, nil
	}
	return &response, nil
}
