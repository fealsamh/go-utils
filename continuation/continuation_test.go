package continuation

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContinuation(t *testing.T) {
	req := require.New(t)

	// unit
	c := Unit[int](1234)
	var s string
	r := c(func(x int) int {
		s = strconv.Itoa(x)
		return x + 1
	})
	req.Equal(1235, r)
	req.Equal("1234", s)

	// bind
	c2 := c.Bind(func(x int) Continuation[float64, int] { return Unit[int](float64(x) + 1) })
	var y float64
	r = c2(func(x float64) int {
		y = x
		return int(x + 1)
	})
	req.Equal(1236, r)
	req.Equal(1235.0, y)

	// fmap
	c3 := c.Fmap(func(x int) float64 { return float64(x) + 2 })
	r = c3(func(x float64) int {
		y = x
		return int(x + 1)
	})
	req.Equal(1237, r)
	req.Equal(1236.0, y)
}
