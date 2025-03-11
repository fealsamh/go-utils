package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// Code is an error code.
type Code int

// error codes
const (
	NotFound Code = iota
	InvalidArgument
	Unauthorised
	Unimplemented
	Internal
)

func (c Code) String() string {
	switch c {
	case NotFound:
		return "not found"
	case InvalidArgument:
		return "invalid argument"
	case Unauthorised:
		return "unauthorised"
	case Unimplemented:
		return "unimplemented"
	case Internal:
		return "internal"
	}
	return ""
}

// HTTPStatus returns the code's HTTP status.
func (c Code) HTTPStatus() int {
	switch c {
	case NotFound:
		return http.StatusNotFound
	case InvalidArgument:
		return http.StatusBadRequest
	case Unauthorised:
		return http.StatusUnauthorized
	case Unimplemented:
		return http.StatusNotImplemented
	case Internal:
		return http.StatusInternalServerError
	}
	return 0
}

// GRPCCode returns the code's gRPC code
func (c Code) GRPCCode() codes.Code {
	switch c {
	case NotFound:
		return codes.NotFound
	case InvalidArgument:
		return codes.InvalidArgument
	case Unauthorised:
		return codes.Unauthenticated
	case Unimplemented:
		return codes.Unimplemented
	case Internal:
		return codes.Internal
	}
	return 0
}
