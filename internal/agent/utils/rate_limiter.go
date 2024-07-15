package utils

type Semaphore struct {
	semaCh chan struct{}
}

// Создание семафора, ограничение нагрухки с помощью максимального количества maxReq.
func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

// Получение задания.
func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

// Завершение задания (высвобоэдение).
func (s *Semaphore) Release() {
	<-s.semaCh
}
