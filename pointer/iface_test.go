package pointer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type iface[T any] interface {
	Pointer[T]
	Dummy() int
}

type impl struct{}

func (x *impl) Dummy() int { return 1234 }

func dummy[T any, PT iface[T]]() int {
	x := PT(new(T))
	return x.Dummy()
}

func TestPointerIface(t *testing.T) {
	req := require.New(t)

	x := dummy[impl]()
	req.Equal(1234, x)
}
