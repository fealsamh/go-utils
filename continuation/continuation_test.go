package continuation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContinuation(t *testing.T) {
	req := require.New(t)

	res, err := Go(func(cont *Continuation[string]) {
		cont.Resume("result", errors.New("error"))
	})

	req.Equal("result", res)
	req.Equal("error", err.Error())
}
