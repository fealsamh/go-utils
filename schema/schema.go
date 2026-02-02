package schema

import (
	"encoding/json"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	schemas = make(map[reflect.Type]*jsonschema.Schema)
)

// Validate validates JSON data against a schema.
func Validate[T any](b []byte) error {
	typ := reflect.TypeFor[T]()
	sch, ok := schemas[typ]
	if !ok {
		var err error
		sch, err = jsonschema.ForType(typ, nil)
		if err != nil {
			return err
		}
		schemas[typ] = sch
	}

	rs, err := sch.Resolve(nil)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	return rs.Validate(m)
}
