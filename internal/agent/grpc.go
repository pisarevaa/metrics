package agent

import (
	"context"
	"errors"

	pb "github.com/pisarevaa/metrics/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Client pb.MetricsClient
	Conn   *grpc.ClientConn
}

// Создание клиента GRPC.
func NewGrpcClient(grpcPort string) (*GrpcClient, error) {
	conn, err := grpc.NewClient(":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	grpcClient := pb.NewMetricsClient(conn)
	return &GrpcClient{Client: grpcClient, Conn: conn}, nil
}

func (c GrpcClient) Close() error {
	err := c.Conn.Close()
	return err
}

func (c GrpcClient) SendMetrics(metrics []Metrics) error {
	var grpcMetrics []*pb.Metric
	for _, metric := range metrics {
		var mType pb.MetricType
		switch metric.MType {
		case gauge:
			mType = pb.MetricType_gauge
		case counter:
			mType = pb.MetricType_counter
		}
		grpcMetrics = append(grpcMetrics, &pb.Metric{
			Id:    metric.ID,
			Type:  mType,
			Value: metric.Value,
			Delta: metric.Delta,
		})
	}
	// добавляем пользователей
	resp, err := c.Client.AddMetrics(context.Background(), &pb.AddMetricsRequest{
		Metrics: grpcMetrics,
	})
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}
	return nil
}
