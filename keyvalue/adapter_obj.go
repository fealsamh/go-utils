package keyvalue

import "reflect"

// ObjectAdapter is a key-value adapter for instances of structures.
type ObjectAdapter[T any] struct {
	obj T
	val reflect.Value
}

// NewObjectAdapter creates a new object adapter.
func NewObjectAdapter[T any](obj T) *ObjectAdapter[T] {
	return &ObjectAdapter[T]{
		obj: obj,
		val: reflect.Indirect(reflect.ValueOf(obj)),
	}
}

// Get returns the value associated with the provided key.
func (a *ObjectAdapter[T]) Get(key string) (interface{}, bool) {
	if v := a.val.FieldByName(key); v.IsValid() {
		return v.Interface(), true
	}
	return nil, false
}

// Put sets a value for the provided key.
func (a *ObjectAdapter[T]) Put(key string, value interface{}) bool {
	if v := a.val.FieldByName(key); v.IsValid() {
		v.Set(reflect.ValueOf(value))
		return true
	}
	return false
}

// Pairs enumerates all the key-value pairs of the underlying instance.
func (a *ObjectAdapter[T]) Pairs(fn func(string, interface{})) {
	t := a.val.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.IsExported() {
			v := a.val.Field(i)
			fn(f.Name, v.Interface())
		}
	}
}

var _ Adapter = new(ObjectAdapter[struct{}])
