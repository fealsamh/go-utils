package equals

import "reflect"

// Deeply deeply compares the two arguments.
func Deeply(x, y any) bool {
	t1, t2 := reflect.TypeOf(x), reflect.TypeOf(y)
	if t1 != t2 {
		return false
	}
	v1, v2 := reflect.ValueOf(x), reflect.ValueOf(y)
	switch t1.Kind() {
	case reflect.Int:
		if v1.Interface().(int) != v2.Interface().(int) {
			return false
		}
	case reflect.Float64:
		if v1.Interface().(float64) != v2.Interface().(float64) {
			return false
		}
	case reflect.String:
		if v1.Interface().(string) != v2.Interface().(string) {
			return false
		}
	case reflect.Pointer:
		if v1.UnsafePointer() != v2.UnsafePointer() {
			switch {
			case v1.IsZero():
				return v2.IsZero()
			case v2.IsZero():
				return v1.IsZero()
			}
			if !Deeply(v1.Elem().Interface(), v2.Elem().Interface()) {
				return false
			}
		}
	case reflect.Struct:
		if !deepEqualStruct(v1.Interface(), v2.Interface()) {
			return false
		}
	}
	return true
}

func deepEqualStruct(x, y any) bool {
	t := reflect.TypeOf(x)
	v1, v2 := reflect.ValueOf(x), reflect.ValueOf(y)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.IsExported() {
			vv1, vv2 := v1.Field(i), v2.Field(i)
			if !Deeply(vv1.Interface(), vv2.Interface()) {
				return false
			}
		}
	}
	return true
}
