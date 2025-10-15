package kvo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObserving(t *testing.T) {
	req := require.New(t)

	obj := NewObject()
	var tuples [][]interface{}
	obj.AddObserver("att1", Prior, func(oldVal, newVal interface{}) {
		val, _ := obj.Get("att1")
		tuples = append(tuples, []interface{}{oldVal, newVal, val})
	})
	obj.AddObserver("att1", Post, func(oldVal, newVal interface{}) {
		val, _ := obj.Get("att1")
		tuples = append(tuples, []interface{}{oldVal, newVal, val})
	})
	obj.Set("att1", 11)
	obj.Set("att1", 22)
	req.Equal([][]interface{}{
		[]interface{}{nil, 11, nil},
		[]interface{}{nil, 11, 11},
		[]interface{}{11, 22, 11},
		[]interface{}{11, 22, 22},
	}, tuples)
}
