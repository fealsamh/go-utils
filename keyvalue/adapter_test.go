package keyvalue

import (
	"reflect"
	"testing"

	"github.com/fealsamh/datastructures/redblack"
)

type s1 struct {
	N  int
	S  *s2
	Sl []*s2
}

type s2 struct {
	T string
}

func TestDeepCopy_ObjToMap(t *testing.T) {
	src, err := NewObjectAdapter(&s1{
		N: 1234,
		S: &s2{
			T: "abcd",
		},
		Sl: []*s2{
			{T: "A"},
			{T: "B"},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	dm := make(map[string]interface{})
	dst := NewMapAdapter(dm)
	if err := DeepCopy(dst, src); err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	m := map[string]interface{}{
		"N": 1234,
		"S": map[string]interface{}{
			"T": "abcd",
		},
		"Sl": []interface{}{
			map[string]interface{}{"T": "A"},
			map[string]interface{}{"T": "B"},
		},
	}
	if !reflect.DeepEqual(m, dm) {
		t.Fatalf("expected '%v', got '%v'", m, dm)
	}
}

func TestDeepCopy_MapToObj(t *testing.T) {
	src := NewMapAdapter(map[string]interface{}{
		"N": 1234,
		"S": map[string]interface{}{
			"T": "abcd",
		},
	})
	ds := new(s1)
	dst, err := NewObjectAdapter(ds)
	if err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	if err := DeepCopy(dst, src); err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	s := &s1{
		N: 1234,
		S: &s2{
			T: "abcd",
		},
	}
	if !reflect.DeepEqual(s, ds) {
		t.Fatalf("expected '%v', got '%v'", s, ds)
	}
}

func TestDeepCopy_ObjToRBTree(t *testing.T) {
	src, err := NewObjectAdapter(&s1{
		N: 1234,
		S: &s2{
			T: "abcd",
		},
		Sl: []*s2{
			{T: "A"},
			{T: "B"},
		},
	})
	if err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	dt := redblack.NewTree[String, interface{}]()
	dst := NewRBTreeAdapter(dt)
	if err := DeepCopy(dst, src); err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	tr := redblack.NewTreeFromMap(map[String]interface{}{
		"N": 1234,
		"S": redblack.NewTreeFromMap(map[String]interface{}{
			"T": "abcd",
		}),
		"Sl": []interface{}{
			redblack.NewTreeFromMap(map[String]interface{}{"T": "A"}),
			redblack.NewTreeFromMap(map[String]interface{}{"T": "B"}),
		},
	})
	if !reflect.DeepEqual(tr, dt) {
		t.Fatalf("expected '%v', got '%v'", tr, dt)
	}
}

func TestDeepCopy_RBTreeToObj(t *testing.T) {
	src := NewRBTreeAdapter(redblack.NewTreeFromMap(map[String]interface{}{
		"N": 1234,
		"S": redblack.NewTreeFromMap(map[String]interface{}{
			"T": "abcd",
		}),
	}))
	ds := new(s1)
	dst, err := NewObjectAdapter(ds)
	if err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	if err := DeepCopy(dst, src); err != nil {
		t.Fatalf("expected nil error, got '%s'", err)
	}
	s := &s1{
		N: 1234,
		S: &s2{
			T: "abcd",
		},
	}
	if !reflect.DeepEqual(s, ds) {
		t.Fatalf("expected '%v', got '%v'", s, ds)
	}
}
