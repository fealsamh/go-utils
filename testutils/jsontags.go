package testutils

import (
	"reflect"
)

// JSONField is a JSON-marshallable field.
type JSONField struct {
	Name string
	Tag  string
}

// JSONFields returns a type's JSON-marshallable fields.
func JSONFields(typ reflect.Type) []JSONField {
	var fields []JSONField
	for _, field := range reflect.VisibleFields(typ) {
		if tag := field.Tag.Get("json"); tag != "-" {
			b := []byte(field.Name)
			fields = append(fields, JSONField{
				Name: string(b),
				Tag:  tag,
			})
		}
	}
	return fields
}
