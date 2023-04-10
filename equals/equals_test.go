package equals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepEqual(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		x, y  any
		equal bool
	}{
		{
			x:     1,
			y:     1,
			equal: true,
		},
		{
			x:     1,
			y:     2,
			equal: false,
		},
		{
			x: struct {
				A int
				B string
			}{A: 1, B: "a"},
			y: struct {
				A int
				B string
			}{A: 1, B: "a"},
			equal: true,
		},
		{
			x: struct {
				A int
				B string
			}{A: 1, B: "a"},
			y: struct {
				A int
				B string
			}{A: 1, B: "b"},
			equal: false,
		},
		{
			x: &struct {
				A int
				B string
			}{A: 1, B: "a"},
			y: &struct {
				A int
				B string
			}{A: 1, B: "a"},
			equal: true,
		},
		{
			x: &struct {
				A int
				B string
			}{A: 1, B: "a"},
			y: &struct {
				A int
				B string
			}{A: 1, B: "b"},
			equal: false,
		},
		{
			x: &struct {
				A *struct{ N int }
			}{},
			y: &struct {
				A *struct{ N int }
			}{A: &struct{ N int }{}},
			equal: false,
		},
	}

	for _, c := range cases {
		assert.Equal(c.equal, Deeply(c.x, c.y))
	}
}
