package pointer

// Pointer is a pointer to [T].
type Pointer[T any] interface {
	*T
}
