package agent

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (ms *MemStorage) Init() {
	ms.Gauge = make(map[string]float64)
	ms.Counter = make(map[string]int64)
}
