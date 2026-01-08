package pointer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {
	req := require.New(t)

	p := To("abcd")
	req.Equal("abcd", *p)
}
