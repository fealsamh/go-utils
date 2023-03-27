package mapreduce

// Engine is a wrapper for map-reduce algorithms.
type Engine[T, V1, V2 any, K1, K2 comparable] struct {
	f1 func(T, func(K1, V1))
	f2 func(K1, []V1) (K2, V2)
}

// Run runs the algorithm.
func (e *Engine[T, V1, V2, K1, K2]) Run(input []T) map[K2]V2 {
	pairs := make(map[K1][]V1)
	for _, x := range input {
		e.f1(x, func(k K1, v V1) {
			r := pairs[k]
			pairs[k] = append(r, v)
		})
	}
	r := make(map[K2]V2)
	for k, l := range pairs {
		k, v := e.f2(k, l)
		r[k] = v
	}
	return r
}

// New creates a new map-reduce engine.
func New[T, V1, V2 any, K1, K2 comparable](f1 func(T, func(K1, V1)), f2 func(K1, []V1) (K2, V2)) *Engine[T, V1, V2, K1, K2] {
	return &Engine[T, V1, V2, K1, K2]{
		f1: f1,
		f2: f2,
	}
}
