package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	config := GetConfigs()
	client := resty.New()
	storage := MemStorage{}
	storage.Init()
	service := Service{Storage: &storage, Client: client, Config: config}
	service.UpdateMetrics()
	heapInuse1 := service.Storage.Gauge["HeapInuse"]
	randomValue1 := service.Storage.Gauge["RandomValue"]
	assert.True(t, len(service.Storage.Gauge) == 28)
	assert.True(t, len(service.Storage.Counter) == 1)
	assert.Equal(t, service.Storage.Counter["PollCount"], int64(1))
	service.UpdateMetrics()
	heapInuse2 := service.Storage.Gauge["HeapInuse"]
	randomValue2 := service.Storage.Gauge["RandomValue"]
	assert.NotEqual(t, heapInuse1, heapInuse2)
	assert.NotEqual(t, randomValue1, randomValue2)
	assert.Equal(t, service.Storage.Counter["PollCount"], int64(2))
}

func TestSendMetrics(t *testing.T) {
	client := resty.New()
	storage := MemStorage{}
	storage.Init()
	service := Service{Storage: &storage, Client: client}
	service.SendMetrics()
	service.UpdateMetrics()
	service.SendMetrics()
	assert.True(t, len(service.Storage.Gauge) == 28)
	assert.True(t, len(service.Storage.Counter) == 1)
}
