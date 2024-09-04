package agent

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/pisarevaa/metrics/internal/agent/utils"
)

type Service struct {
	Client     *resty.Client      // клиент для внешних запросов
	Storage    *MemStorage        // хранилище метрик
	Config     Config             // параметры конфигурации
	Logger     *zap.SugaredLogger // логер
	Semaphore  *utils.Semaphore   // семафор
	GrpcClient *GrpcClient        // GRPC client
}

// Созание нового сервиса.
func NewService(
	client *resty.Client,
	storage *MemStorage,
	config Config,
	logger *zap.SugaredLogger,
	semaphore *utils.Semaphore,
	grpcClient *GrpcClient,
) *Service {
	return &Service{
		Client:     client,
		Storage:    storage,
		Config:     config,
		Logger:     logger,
		Semaphore:  semaphore,
		GrpcClient: grpcClient,
	}
}
