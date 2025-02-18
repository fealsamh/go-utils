package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error is an error with an error code.
type Error struct {
	message string
	code    Code
}

func (e *Error) Error() string {
	return e.message
}

// New creates a new error.
func New(message string, code Code) *Error {
	return &Error{message, code}
}

// WriteHTTPHeader writes the HTTP header.
func (e *Error) WriteHTTPHeader(w http.ResponseWriter) {
	w.WriteHeader(e.code.HTTPStatus())
}

// GRPCError returns the corresponding GRPCError.
func (e *Error) GRPCError() error {
	return status.Error(e.code.GRPCCode(), e.message)
}

// ConvertToGRPC converts the error into a gRPC error.
func ConvertToGRPC(err error) error {
	if err, ok := err.(*Error); ok {
		return err.GRPCError()
	}
	return status.Error(codes.Internal, err.Error())
}

// WriteHTTPHeader writes the HTTP header.
func WriteHTTPHeader(err error, w http.ResponseWriter) {
	if err, ok := err.(*Error); ok {
		err.WriteHTTPHeader(w)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
