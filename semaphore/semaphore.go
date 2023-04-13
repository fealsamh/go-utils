package semaphore

// Semaphore is a semaphore for synchronisation.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with an initial value.
func NewSemaphore(val int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, val),
	}
}

// Wait decreases the semaphore value, waiting if it gets beneath zero.
func (s *Semaphore) Wait() {
	s.ch <- struct{}{}
}

// Signal increases the semaphore value.
func (s *Semaphore) Signal() {
	<-s.ch
}

