package errors

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

func TestUUIDError(t *testing.T) {
	req := require.New(t)

	_, err := uuid.Parse("x")
	req.Equal(invalidUUIDErrMessage, err.Error())
}
