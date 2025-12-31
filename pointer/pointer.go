package pointer

type simpleType interface {
	int | uint | int32 | uint32 | int64 | uint64 | bool | float32 | float64 | string
}

// To returns a pointer to its argument.
func To[T simpleType](x T) *T {
	return &x
}

// Pointer is a pointer to [T].
type Pointer[T any] interface {
	*T
}
