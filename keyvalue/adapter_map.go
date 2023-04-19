package keyvalue

import "reflect"

// MapAdapter is a key-value adapter for in-built hash maps keyed off of strings.
type MapAdapter struct {
	m map[string]interface{}
}

// NewMapAdapter creates a new map adapter.
func NewMapAdapter(m map[string]interface{}) *MapAdapter {
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

// ShouldConvert returns true if the provided type should be converted in the course of copying.
func (a *MapAdapter) ShouldConvert(t reflect.Type) bool {
	return t == stringAnyMapType
}

// NewInstance creates a new instance of the provided type and an associated adapter.
func (a *MapAdapter) NewInstance(reflect.Type) (interface{}, Adapter, error) {
	m := make(map[string]interface{})
	return m, NewMapAdapter(m), nil
}

// NewAdapter creates a new key-value adapter for the provided value.
func (a *MapAdapter) NewAdapter(v interface{}) (Adapter, error) {
	return NewMapAdapter(v.(map[string]interface{})), nil
}

// TypeForKey returns the type of values associated with the key.
func (a *MapAdapter) TypeForKey(key string) reflect.Type {
	return stringAnyMapType
}

var (
	stringAnyMapType = reflect.TypeOf((map[string]interface{})(nil))

	_ Adapter = new(MapAdapter)
)
