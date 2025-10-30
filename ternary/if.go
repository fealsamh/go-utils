package ternary

// If is the usual ternary if statement.
// Note that all arguments are evaluated upon call.
func If[T any](cond bool, ifTrue, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}
