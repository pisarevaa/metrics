package agent

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Service struct {
	Client    *resty.Client
	Storage   *MemStorage
	Config    Config
	Logger    *zap.SugaredLogger
	Semaphore *Semaphore
}

func NewService(
	client *resty.Client,
	storage *MemStorage,
	config Config,
	logger *zap.SugaredLogger,
	semaphore *Semaphore,
) *Service {
	return &Service{
		Client:    client,
		Storage:   storage,
		Config:    config,
		Logger:    logger,
		Semaphore: semaphore,
	}
}
