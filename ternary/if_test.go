package ternary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIf(t *testing.T) {
	require := require.New(t)

	r1 := If(true, "a", "b")
	require.Equal("a", r1)

	r2 := If(false, 1, 2)
	require.Equal(2, r2)
}
