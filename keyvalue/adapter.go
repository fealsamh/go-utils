package keyvalue

// Pair is a key-value pair.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// Adapter is a key-value adapter.
type Adapter interface {
	Get(string) (interface{}, bool)
	Put(string, interface{})
	Pairs(func(string, interface{}))
}
