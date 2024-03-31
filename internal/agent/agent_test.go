package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	storage := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	storage.UpdateMetrics()
	heapInuse1 := storage.Gauge["HeapInuse"]
	randomValue1 := storage.Gauge["RandomValue"]
	assert.True(t, len(storage.Gauge) == 28)
	assert.True(t, len(storage.Counter) == 1)
	assert.Equal(t, storage.Counter["PollCount"], int64(1))
	storage.UpdateMetrics()
	heapInuse2 := storage.Gauge["HeapInuse"]
	randomValue2 := storage.Gauge["RandomValue"]
	assert.NotEqual(t, heapInuse1, heapInuse2)
	assert.NotEqual(t, randomValue1, randomValue2)
	assert.Equal(t, storage.Counter["PollCount"], int64(2))
}

func TestSendMetrics(t *testing.T) {
	storage := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	storage.SendMetrics()
	storage.UpdateMetrics()
	storage.SendMetrics()
	assert.True(t, len(storage.Gauge) == 28)
	assert.True(t, len(storage.Counter) == 1)
}
