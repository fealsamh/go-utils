package schema

import (
	"testing"

	"github.com/fealsamh/go-utils/nocopy"
	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	req := require.New(t)

	type s struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := Validate[s](nocopy.Bytes(`{"name":"Jane", "age": 18}`))
	req.Nil(err)
}
