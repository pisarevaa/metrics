package agent_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/pisarevaa/metrics/internal/agent"
)

func TestUpdateMetrics(t *testing.T) {
	config := agent.GetConfig()
	client := resty.New()
	storage := agent.NewMemStorageRepo()
	service := agent.Service{Storage: storage, Client: client, Config: config}
	service.UpdateMetrics()
	heapInuseFirst, heapInuseFirstErr := service.Storage.Get("gauge", "HeapInuse")
	if assert.NoError(t, heapInuseFirstErr) {
		assert.NotEmpty(t, heapInuseFirst)
	}
	randomValueFirst, randomValueFirstErr := service.Storage.Get("gauge", "RandomValue")
	if assert.NoError(t, randomValueFirstErr) {
		assert.NotEmpty(t, randomValueFirst)
	}
	pollCounterFirst, pollCounterFirstErr := service.Storage.Get("counter", "PollCount")
	if assert.NoError(t, pollCounterFirstErr) {
		assert.Equal(t, "1", pollCounterFirst)
	}
	service.UpdateMetrics()
	heapInuseSecond, heapInuseSecondErr := service.Storage.Get("gauge", "HeapInuse")
	if assert.NoError(t, heapInuseSecondErr) {
		assert.NotEmpty(t, heapInuseSecond)
		assert.NotEqual(t, heapInuseSecond, heapInuseFirst)
	}
	randomValueSecond, randomValueSecondErr := service.Storage.Get("gauge", "RandomValue")
	if assert.NoError(t, randomValueSecondErr) {
		assert.NotEmpty(t, randomValueSecond)
		assert.NotEqual(t, randomValueSecond, randomValueFirst)
	}
	pollCounterSecond, pollCounterSecondErr := service.Storage.Get("counter", "PollCount")
	if assert.NoError(t, pollCounterSecondErr) {
		assert.Equal(t, "2", pollCounterSecond)
	}
}

func TestSendMetrics(t *testing.T) {
	client := resty.New()
	storage := agent.NewMemStorageRepo()
	service := agent.Service{Storage: storage, Client: client}
	service.SendMetrics()
	service.UpdateMetrics()
	service.SendMetrics()
	assert.NotEmpty(t, service.Storage.GetAll())
}
