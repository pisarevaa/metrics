package agent_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/pisarevaa/metrics/internal/agent"
)

func TestUpdateMetrics(t *testing.T) {
	config := agent.GetConfigs()
	client := resty.New()
	storage := agent.MemStorage{}
	storage.Init()
	service := agent.Service{Storage: &storage, Client: client, Config: config}
	service.UpdateMetrics()
	heapInuse1 := service.Storage.Gauge["HeapInuse"]
	randomValue1 := service.Storage.Gauge["RandomValue"]
	assert.Len(t, service.Storage.Gauge, 28)
	assert.Len(t, service.Storage.Counter, 1)
	assert.Equal(t, int64(1), service.Storage.Counter["PollCount"])
	service.UpdateMetrics()
	heapInuse2 := service.Storage.Gauge["HeapInuse"]
	randomValue2 := service.Storage.Gauge["RandomValue"]
	assert.NotEqual(t, heapInuse1, heapInuse2)
	assert.NotEqual(t, randomValue1, randomValue2)
	assert.Equal(t, int64(2), service.Storage.Counter["PollCount"])
}

func TestSendMetrics(t *testing.T) {
	client := resty.New()
	storage := agent.MemStorage{}
	storage.Init()
	service := agent.Service{Storage: &storage, Client: client}
	service.SendMetrics()
	service.UpdateMetrics()
	service.SendMetrics()
	assert.Len(t, service.Storage.Gauge, 28)
	assert.Len(t, service.Storage.Counter, 1)
}
