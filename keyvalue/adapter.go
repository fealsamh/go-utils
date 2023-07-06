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
	TypeForCompoundKey(string) reflect.Type
	TypeForSliceKey(string) reflect.Type
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
		var v interface{}
		v, err = deepCopyValue(value, dst.TypeForCompoundKey(key), dst.TypeForSliceKey(key), src, dst)
		if err == nil {
			err = dst.Put(key, v)
		}
		return err == nil
	})
	return err
}

func deepCopyValue(value interface{}, ty, tys reflect.Type, src, dst Adapter) (interface{}, error) {
	st := reflect.TypeOf(value)
	if st.Kind() == reflect.Slice {
		sv := reflect.ValueOf(value)
		if sv.IsZero() {
			return nil, nil
		}
		ds := reflect.MakeSlice(tys, 0, sv.Len())
		for i := 0; i < sv.Len(); i++ {
			el, err := deepCopyValue(sv.Index(i).Interface(), ty, nil, src, dst)
			if err != nil {
				return nil, err
			}
			ds = reflect.Append(ds, reflect.ValueOf(el))
		}
		return ds.Interface(), nil
	}
	if src.ShouldConvert(st) {
		sa, err := src.NewAdapter(value)
		if err != nil {
			return nil, err
		}
		dt := ty
		v, da, err := dst.NewInstance(dt)
		if err != nil {
			return nil, err
		}
		if err := DeepCopy(da, sa); err != nil {
			return nil, err
		}
		return v, nil
	}
	return value, nil
}
