package keyvalue

import (
	"fmt"
	"reflect"
)

// ObjectAdapter is a key-value adapter for instances of structures.
type ObjectAdapter[T any] struct {
	obj T
	val reflect.Value
}

// NewObjectAdapter creates a new object adapter.
func NewObjectAdapter[T any](obj T) (*ObjectAdapter[T], error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Pointer && val.IsZero() {
		return nil, fmt.Errorf("failed to create object adapter for nil pointer (%s)", val.Type())
	}
	return &ObjectAdapter[T]{
		obj: obj,
		val: reflect.Indirect(val),
	}, nil
}

// Get returns the value associated with the provided key.
func (a *ObjectAdapter[T]) Get(key string) (interface{}, bool) {
	if v := a.val.FieldByName(key); v.IsValid() {
		return v.Interface(), true
	}
	return nil, false
}

// Put sets a value for the provided key.
func (a *ObjectAdapter[T]) Put(key string, value interface{}) error {
	if v := a.val.FieldByName(key); v.IsValid() {
		vs := reflect.ValueOf(value)
		if !vs.Type().ConvertibleTo(v.Type()) {
			return fmt.Errorf("'put' for key '%s' failed, cannot convert '%s' into '%s'", key, vs.Type().Name(), v.Type().Name())
		}
		v.Set(vs.Convert(v.Type()))
		return nil
	}
	return fmt.Errorf("'put' for key '%s' failed, no such key in type", key)
}

// Pairs enumerates all the key-value pairs of the underlying instance.
func (a *ObjectAdapter[T]) Pairs(fn func(string, interface{}) bool) bool {
	t := a.val.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.IsExported() {
			v := a.val.Field(i)
			if !fn(f.Name, v.Interface()) {
				return false
			}
		}
	}
	return true
}

var _ Adapter = new(ObjectAdapter[struct{}])
