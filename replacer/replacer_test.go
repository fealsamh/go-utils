package replacer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplacer(t *testing.T) {
	req := require.New(t)

	m := map[string]string{
		"a": "A",
		"b": "B",
	}

	req.Equal("cac", Replace("{c}a{c}", m))
	req.Equal("AaB", Replace("{a}a{b}", m))
	req.Equal("1AB2", Replace("1{a}{b}2", m))
}
