package keyvalue

// MapAdapter is a key-value adapter for in-built hash maps keyed off of strings.
type MapAdapter struct {
	Map map[string]interface{}
}

// NewMapAdapter creates a new map adapter.
func NewMapAdapter(m map[string]interface{}) *MapAdapter {
	if m == nil {
		m = make(map[string]interface{})
	}
	return &MapAdapter{Map: m}
}

// Get returns the value associated with the provided key.
func (a *MapAdapter) Get(key string) (interface{}, bool) {
	value, ok := a.Map[key]
	return value, ok
}

// Put sets a value for the provided key.
func (a *MapAdapter) Put(key string, value interface{}) {
	a.Map[key] = value
}

// Pairs enumerate all key-value pairs of the underlying map.
func (a *MapAdapter) Pairs(f func(string, interface{})) {
	for k, v := range a.Map {
		f(k, v)
	}
}

var _ Adapter = new(MapAdapter)
