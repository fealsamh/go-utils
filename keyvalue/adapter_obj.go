package keyvalue

import (
	"fmt"
	"reflect"
)

// ObjectAdapter is a key-value adapter for instances of structures.
type ObjectAdapter struct {
	val reflect.Value
}

// NewObjectAdapter creates a new object adapter.
func NewObjectAdapter[T any](obj T) (*ObjectAdapter, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("failed to create object adapter for non-pointer (%s)", val.Type())
	}
	if val.IsZero() {
		return nil, fmt.Errorf("failed to create object adapter for nil pointer (%s)", val.Type())
	}
	return &ObjectAdapter{
		val: reflect.Indirect(val),
	}, nil
}

// Get returns the value associated with the provided key.
func (a *ObjectAdapter) Get(key string) (interface{}, bool) {
	if v := a.val.FieldByName(key); v.IsValid() {
		return v.Interface(), true
	}
	return nil, false
}

// Put sets a value for the provided key.
func (a *ObjectAdapter) Put(key string, value interface{}) error {
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
func (a *ObjectAdapter) Pairs(fn func(string, interface{}) bool) bool {
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

// ShouldConvert returns true if the provided type should be converted in the course of copying.
func (a *ObjectAdapter) ShouldConvert(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}

// NewInstance creates a new instance of the provided type and an associated adapter.
func (a *ObjectAdapter) NewInstance(t reflect.Type) (interface{}, Adapter, error) {
	v := reflect.New(t.Elem()).Interface()
	a, err := NewObjectAdapter(v)
	if err != nil {
		return nil, nil, err
	}
	return v, a, nil
}

// NewAdapter creates a new key-value adapter for the provided value.
func (a *ObjectAdapter) NewAdapter(v interface{}) (Adapter, error) {
	return NewObjectAdapter(v)
}

// TypeForKey returns the type of values associated with the key.
func (a *ObjectAdapter) TypeForKey(key string) reflect.Type {
	f, ok := a.val.Type().FieldByName(key)
	if !ok {
		panic(fmt.Sprintf("unknown key '%s' in type '%s'", key, a.val.Type()))
	}
	return f.Type
}

var _ Adapter = new(ObjectAdapter)
