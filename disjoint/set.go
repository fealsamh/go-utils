package disjoint

// Ordered is an ordered type.
type Ordered[T any] interface {
	Less(T) bool
}

// SetEl is a disjoint set's elements.
type SetEl[T any] struct {
	Value  T
	parent *SetEl[T]
}

// NewEl creates a new element of a disjoint set.
func NewEl[T any](x T) *SetEl[T] {
	return &SetEl[T]{Value: x}
}

func (x *SetEl[T]) Find() *SetEl[T] {
	if x.parent == nil {
		return x
	}
	return x.parent.Find()
}

func (x *SetEl[T]) Union(y *SetEl[T]) {
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
