package schema

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	// schemas = make(map[reflect.Type]*jsonschema.Resolved)
	schemas sync.Map
)

// Validate validates JSON data against a schema.
func Validate[T any](b []byte) error {
	typ := reflect.TypeFor[T]()
	rs, ok := schemas.Load(typ)
	if !ok {
		sch, err := jsonschema.ForType(typ, nil)
		if err != nil {
			return err
		}
		rs, err = sch.Resolve(nil)
		if err != nil {
			return err
		}
		schemas.Store(typ, rs)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	return rs.(*jsonschema.Resolved).Validate(m)
}
