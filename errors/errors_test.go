package errors

import (
	"errors"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

func TestFromErrorInvalidUUID(t *testing.T) {
	req := require.New(t)

	_, stdlibErr := uuid.Parse("nope")
	req.EqualError(stdlibErr, "invalid uuid")

	for _, tc := range []struct {
		name string
		err  error
		code Code
		ok   bool
	}{
		{name: "stdlib parse failure", err: stdlibErr, code: InvalidArgument, ok: true},
		{name: "google invalid length", err: errors.New("invalid UUID length: 4"), code: InvalidArgument, ok: true},
		{name: "google invalid format", err: errors.New("invalid UUID format"), code: InvalidArgument, ok: true},
		{name: "unrelated error", err: errors.New("boom"), ok: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			wrapped, ok := FromError(tc.err)
			req.Equal(tc.ok, ok)
			if tc.ok {
				req.Equal(tc.code, wrapped.code)
			}
		})
	}
}
