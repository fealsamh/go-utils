package errors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
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

// WrappedError is a wrapped error with an error code.
type WrappedError struct {
	err  error
	code Code
}

func (e *WrappedError) Error() string {
	return e.err.Error()
}

func (e *WrappedError) Unwrap() error {
	return e.err
}

// New creates a new error.
func New(message string, code Code) *Error {
	return &Error{message, code}
}

// Wrap creates a wrapped error.
func Wrap(err error, code Code) *WrappedError {
	return &WrappedError{err, code}
}

// FromError creates an error with an error code from the provided error.
func FromError(err error) (*WrappedError, bool) {
	switch {
	case errors.Is(err, errors.ErrUnsupported):
		return &WrappedError{err, Unimplemented}, true

	case errors.Is(err, sql.ErrNoRows):
		return &WrappedError{err, NotFound}, true

	case uuid.IsInvalidLengthError(err):
		return &WrappedError{err, InvalidArgument}, true

	case err.Error() == "invalid UUID format":
		return &WrappedError{err, InvalidArgument}, true
	}

	var jsonErr *json.SyntaxError
	if errors.As(err, &jsonErr) {
		return &WrappedError{err, InvalidArgument}, true
	}

	return nil, false
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
	if err := convertToGRPC(err); err != nil {
		return err
	}
	return status.Error(codes.Internal, err.Error())
}

func convertToGRPC(err error) error {
	if err, ok := err.(*Error); ok {
		return err.GRPCError()
	}
	if err, ok := err.(wrappedError); ok {
		return convertToGRPC(err.Unwrap())
	}
	if err, ok := err.(wrappedErrors); ok {
		for _, err := range err.Unwrap() {
			if err := convertToGRPC(err); err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteHTTPError writes the HTTP error.
func WriteHTTPError(w http.ResponseWriter, err error) {
	WriteHTTPHeader(w, err)
	io.WriteString(w, err.Error())
}

type jsonError struct {
	Error string `json:"error"`
}

// WriteHTTPErrorJSON writes the HTTP error as JSON.
func WriteHTTPErrorJSON(w http.ResponseWriter, err error) {
	WriteHTTPHeader(w, err)
	json.NewEncoder(w).Encode(jsonError{err.Error()})
}

// WriteHTTPHeader writes the HTTP header.
func WriteHTTPHeader(w http.ResponseWriter, err error) {
	if s, ok := httpStatus(err); ok {
		w.WriteHeader(s)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func httpStatus(err error) (int, bool) {
	if err, ok := err.(*Error); ok {
		return err.code.HTTPStatus(), true
	}
	if err, ok := err.(wrappedError); ok {
		return httpStatus(err.Unwrap())
	}
	if err, ok := err.(wrappedErrors); ok {
		for _, err := range err.Unwrap() {
			if s, ok := httpStatus(err); ok {
				return s, true
			}
		}
	}
	return 0, false
}
