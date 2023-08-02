package ternary

// If is the usual ternary if statement.
func If[T any](cond bool, ifTrue, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}
