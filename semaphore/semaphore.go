package semaphore

import "context"

// Semaphore is a semaphore for synchronisation.
type Semaphore struct {
	ch chan struct{}
}

// New creates a new semaphore with an initial value.
func New(val int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, val),
	}
}

// Wait decreases the semaphore value, waiting if it gets beneath zero.
func (s *Semaphore) Wait() {
	s.ch <- struct{}{}
}

// WaitCtx decreases the semaphore value, waiting if it gets beneath zero.
// The provided context can be used to stop waiting.
func (s *Semaphore) WaitCtx(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Signal increases the semaphore value.
func (s *Semaphore) Signal() {
	<-s.ch
}
