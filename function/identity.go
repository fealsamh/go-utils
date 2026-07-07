package function

// Identity is the identity function.
func Identity[T any](x T) T {
	return x
}
