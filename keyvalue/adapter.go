package keyvalue

import (
	"reflect"
)

// Adapter is a key-value adapter.
type Adapter interface {
	Get(string) (interface{}, bool)
	Put(string, interface{}) error
	Pairs(func(string, interface{}) bool) bool
	ShouldConvert(reflect.Type) bool
	NewInstance(reflect.Type) (interface{}, Adapter, error)
	NewAdapter(interface{}) (Adapter, error)
	TypeForKey(string) reflect.Type
}

// Copy copies the contents of `src` to `dst`.
func Copy(dst, src Adapter) error {
	var err error
	src.Pairs(func(key string, value interface{}) bool {
		err = dst.Put(key, value)
		return err == nil
	})
	return err
}

// DeepCopy deeply copies the contents of `src` to `dst`.
func DeepCopy(dst, src Adapter) error {
	var err error
	src.Pairs(func(key string, value interface{}) bool {
		st := reflect.TypeOf(value)
		if src.ShouldConvert(st) {
			var (
				v  interface{}
				sa Adapter
				da Adapter
			)
			if sa, err = src.NewAdapter(value); err == nil {
				dt := dst.TypeForKey(key)
				if v, da, err = dst.NewInstance(dt); err == nil {
					if err = DeepCopy(da, sa); err == nil {
						err = dst.Put(key, v)
					}
				}
			}
		} else {
			err = dst.Put(key, value)
		}
		return err == nil
	})
	return err
}
