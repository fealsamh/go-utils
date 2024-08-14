package disjoint

type Ordered[T any] interface {
	Less(T) bool
}

type Set[T any] struct {
	Value  T
	parent *Set[T]
}

func New[T any](x T) *Set[T] {
	return &Set[T]{Value: x}
}

func (x *Set[T]) Find() *Set[T] {
	if x.parent == nil {
		return x
	}
	return x.parent.Find()
}

func (x *Set[T]) Union(y *Set[T]) {
	v1 := x.Find()
	v2 := y.Find()
	if v, ok := any(v1.Value).(Ordered[T]); ok {
		if v.Less(v2.Value) {
			v2.parent = v1
		} else {
			v1.parent = v2
		}
	} else {
		v1.parent = v2
	}
}
