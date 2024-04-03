package main

import (
	"fmt"
	"github.com/pisarevaa/metrics/internal/agent"
	"github.com/go-resty/resty/v2"
	"sync"
)

func main() {
	config := agent.GetConfigs()
	client := resty.New()
	storage := agent.MemStorage{}
	storage.Init()
	service := agent.Service{Storage: &storage, Client: client, Config: config}
	var wg sync.WaitGroup
	wg.Add(2)
	fmt.Println("Client is running...")
	go service.RunUpdateMetrics(&wg)
	go service.RunSendMetrics(&wg)
	wg.Wait()
}
