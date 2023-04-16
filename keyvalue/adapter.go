package keyvalue

// Pair is a key-value pair.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// Adapter is a key-value adapter.
type Adapter interface {
	Get(string) (interface{}, bool)
	Put(string, interface{}) bool
	Pairs(func(string, interface{}))
}

// Copy copies the contents of `src` to `dst`.
func Copy(dst, src Adapter) bool {
	var failed bool
	src.Pairs(func(key string, value interface{}) {
		if !dst.Put(key, value) {
			failed = true
		}
	})
	return !failed
}
