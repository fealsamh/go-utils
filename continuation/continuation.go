package continuation

import "sync"

// Continuation is a continuation.
type Continuation[T any] struct {
	wg     *sync.WaitGroup
	result T
	err    error
}

// Resume resumes the continuation.
func (c *Continuation[T]) Resume(result T, err error) {
	c.result = result
	c.err = err
	c.wg.Done()
}

// Go creates a continuation.
func Go[T any](f func(*Continuation[T])) (T, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	cont := &Continuation[T]{
		wg: &wg,
	}
	go f(cont)
	wg.Wait()
	return cont.result, cont.err
}
