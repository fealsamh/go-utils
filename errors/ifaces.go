package errors

type wrappedError interface {
	Unwrap() error
}

type wrappedErrors interface {
	Unwrap() []error
}
