package errors

import (
	"encoding/json"
	"io"
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
func WriteHTTPError(err error, w http.ResponseWriter) {
	WriteHTTPHeader(err, w)
	io.WriteString(w, err.Error())
}

type jsonError struct {
	Error string `json:"error"`
}

// WriteHTTPErrorJSON writes the HTTP error as JSON.
func WriteHTTPErrorJSON(err error, w http.ResponseWriter) {
	WriteHTTPHeader(err, w)
	json.NewEncoder(w).Encode(jsonError{err.Error()})
}

// WriteHTTPHeader writes the HTTP header.
func WriteHTTPHeader(err error, w http.ResponseWriter) {
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
