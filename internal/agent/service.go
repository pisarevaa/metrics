package agent

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/pisarevaa/metrics/internal/agent/utils"
)

type Service struct {
	Client    *resty.Client
	Storage   *MemStorage
	Config    Config
	Logger    *zap.SugaredLogger
	Semaphore *utils.Semaphore
}

func NewService(
	client *resty.Client,
	storage *MemStorage,
	config Config,
	logger *zap.SugaredLogger,
	semaphore *utils.Semaphore,
) *Service {
	return &Service{
		Client:    client,
		Storage:   storage,
		Config:    config,
		Logger:    logger,
		Semaphore: semaphore,
	}
}
