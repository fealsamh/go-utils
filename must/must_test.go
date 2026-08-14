package must

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func testFunc(x int, succeed bool) (int, error) {
	if succeed {
		return x, nil
	}
	return 0, errors.New("some error")
}

func TestMustSuccess(t *testing.T) {
	req := require.New(t)

	req.Equal(12, Must(testFunc(12, true)))
}

func TestMustFailure(t *testing.T) {
	req := require.New(t)

	req.Panics(func() { Must(testFunc(12, false)) })
}
