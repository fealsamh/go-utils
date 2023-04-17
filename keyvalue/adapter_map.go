package keyvalue

// MapAdapter is a key-value adapter for in-built hash maps keyed off of strings.
type MapAdapter struct {
	m map[string]interface{}
}

// NewMapAdapter creates a new map adapter.
func NewMapAdapter(m map[string]interface{}) *MapAdapter {
	if m == nil {
		m = make(map[string]interface{})
	}
	return &MapAdapter{m: m}
}

// Get returns the value associated with the provided key.
func (a *MapAdapter) Get(key string) (interface{}, bool) {
	value, ok := a.m[key]
	return value, ok
}

// Put sets a value for the provided key.
func (a *MapAdapter) Put(key string, value interface{}) error {
	a.m[key] = value
	return nil
}

// Pairs enumerates all the key-value pairs of the underlying map.
func (a *MapAdapter) Pairs(f func(string, interface{}) bool) bool {
	for k, v := range a.m {
		if !f(k, v) {
			return false
		}
	}
	return true
}

var _ Adapter = new(MapAdapter)
