package keyvalue

import (
	"reflect"
	"strings"

	"github.com/fealsamh/datastructures/redblack"
)

// String is a totally-comparable string.
type String string

// Compare lexicographically compares two strings.
func (s String) Compare(s2 String) int { return strings.Compare(string(s), string(s2)) }

// RBTreeAdapter is a key-value adapter for red-black trees keyed off of strings.
type RBTreeAdapter struct {
	tr *redblack.Tree[String, interface{}]
}

// NewRBTreeAdapter creates a new red-black tree adapter.
func NewRBTreeAdapter(tr *redblack.Tree[String, interface{}]) *RBTreeAdapter {
	return &RBTreeAdapter{tr: tr}
}

// Get returns the value associated with the provided key.
func (a *RBTreeAdapter) Get(key string) (interface{}, bool) {
	return a.tr.Get(String(key))
}

// Put sets a value for the provided key.
func (a *RBTreeAdapter) Put(key string, value interface{}) error {
	a.tr.Put(String(key), value)
	return nil
}

// Pairs enumerates all the key-value pairs of the underlying map.
func (a *RBTreeAdapter) Pairs(f func(string, interface{}) bool) bool {
	return a.tr.Enumerate(func(k String, v interface{}) bool {
		return f(string(k), v)
	})
}

// ShouldConvert returns true if the provided type should be converted in the course of copying.
func (a *RBTreeAdapter) ShouldConvert(t reflect.Type) bool {
	return t == stringAnyRBTreeType
}

// NewInstance creates a new instance of the provided type and an associated adapter.
func (a *RBTreeAdapter) NewInstance(reflect.Type) (interface{}, Adapter, error) {
	tr := redblack.NewTree[String, interface{}]()
	return tr, NewRBTreeAdapter(tr), nil
}

// NewAdapter creates a new key-value adapter for the provided value.
func (a *RBTreeAdapter) NewAdapter(v interface{}) (Adapter, error) {
	return NewRBTreeAdapter(v.(*redblack.Tree[String, interface{}])), nil
}

// TypeForCompoundKey returns the type of values associated with the key.
func (a *RBTreeAdapter) TypeForCompoundKey(key string) reflect.Type {
	return stringAnyRBTreeType
}

// TypeForSliceKey returns the type of values associated with the key which is a slice.
func (a *RBTreeAdapter) TypeForSliceKey(key string) reflect.Type {
	return emptyInterfaceSliceType
}

var (
	stringAnyRBTreeType = reflect.TypeOf((*redblack.Tree[String, interface{}])(nil))

	_ Adapter = new(RBTreeAdapter)
)
