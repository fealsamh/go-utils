package kvo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObserving(t *testing.T) {
	req := require.New(t)

	obj := NewObject()
	var tuples [][]any
	obj.AddObserver("att1", Prior, func(oldVal, newVal any) {
		val, _ := obj.Get("att1")
		tuples = append(tuples, []any{oldVal, newVal, val})
	})
	obj.AddObserver("att1", Post, func(oldVal, newVal any) {
		val, _ := obj.Get("att1")
		tuples = append(tuples, []any{oldVal, newVal, val})
	})
	obj.Set("att1", 11)
	obj.Set("att1", 22)
	req.Equal([][]any{
		[]any{nil, 11, nil},
		[]any{nil, 11, 11},
		[]any{11, 22, 11},
		[]any{11, 22, 22},
	}, tuples)
}
