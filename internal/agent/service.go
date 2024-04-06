package agent

import (
	"github.com/go-resty/resty/v2"
)

type Service struct {
	Client  *resty.Client
	Storage *MemStorage
	Config  Config
}

func NewService(client *resty.Client, storage *MemStorage, config Config) *Service {
	return &Service{
		Client:  client,
		Storage: storage,
		Config:  config,
	}
}
