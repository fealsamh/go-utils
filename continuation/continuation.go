package continuation

import "github.com/fealsamh/go-utils/function"

// Continuation is a continuation.
type Continuation[T, R any] func(func(T) R) R

// Unit ...
func Unit[R, T any](x T) Continuation[T, R] {
	return func(f func(T) R) R {
		return f(x)
	}
}

// Bind ...
func (c Continuation[T, R]) Bind[U any](f func(T) Continuation[U, R]) Continuation[U, R] {
	return func(g func(U) R) R {
		return c(func(x T) R {
			return f(x)(g)
		})
	}
}

// Fmap ...
func (c Continuation[T, R]) Fmap[U any](f func(T) U) Continuation[U, R] {
	return c.Bind(func(x T) Continuation[U, R] {
		return Unit[R](f(x))
	})
}

// Join ...
func Join[T, R any](c Continuation[Continuation[T, R], R]) Continuation[T, R] {
	return c.Bind(function.Identity)
}
