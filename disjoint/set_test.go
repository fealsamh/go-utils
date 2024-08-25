package disjoint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type String string

func (s1 String) Less(s2 String) bool {
	return s1 < s2
}

func TestUnion(t *testing.T) {
	req := require.New(t)

	x := NewEl[String]("a")
	y := NewEl[String]("b")

	req.Equal("a", string(x.Find().Value))
	req.Equal("b", string(y.Find().Value))

	x.Union(y)
	req.Equal("a", string(x.Find().Value))
	req.Equal("a", string(y.Find().Value))
}
