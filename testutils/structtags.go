package testutils

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/fealsamh/go-utils/textutils"
)

// Case is a specific text case.
type Case int

const (
	// SnakeCase ...
	SnakeCase Case = iota
	// CamelCase ...
	CamelCase
)

// CheckStructTags checks whether field names match with struct tags.
func CheckStructTags(tagAttr string, cas Case, typ reflect.Type, terms ...string) error {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if slices.Index([]string{"0", "no", "false"}, f.Tag.Get("check")) != -1 {
			continue
		}
		tag := f.Tag.Get(tagAttr)
		if tag == "" {
			if f.Anonymous {
				continue
			}
			return fmt.Errorf("nil '%s' struct tag for %s.%s", tagAttr, typ, f.Name)
		}
		if !f.Anonymous {
			tagName := strings.Split(tag, ",")[0]
			if f.Name == "XMLName" {
				typName := convertCase(cas, typ.Name())
				comps := strings.Split(tagName, " ")
				tagName = comps[len(comps)-1]
				if strings.ToLower(typName) != tagName {
					return fmt.Errorf("bad '%s' struct tag for %s.XMLName (%s != %s)", tagAttr, typ, typName, tagName)
				}
			} else if tagName != "-" {
				snakeField := convertCase(cas, f.Name, terms...)
				if snakeField != tagName {
					return fmt.Errorf("bad '%s' struct tag for %s.%s (%s != %s)", tagAttr, typ, f.Name, snakeField, tagName)
				}
			}
		}
	}
	return nil
}

func convertCase(cas Case, s string, terms ...string) string {
	switch cas {
	case SnakeCase:
		return textutils.SnakecasedFromCamelcased(s, terms...)
	case CamelCase:
		comps := textutils.SplitCamelcasedString(s, terms...)
		var s string
		for i, c := range comps {
			lc := strings.ToLower(c)
			if slices.Index(terms, lc) != -1 {
				if i == 0 {
					s += lc
				} else {
					s += strings.ToUpper(lc[:1]) + lc[1:]
				}
			} else {
				if i == 0 {
					s += strings.ToLower(c[:1]) + c[1:]
				} else {
					s += c
				}
			}
		}
		return s
	default:
		panic("unknown case")
	}
}

// CheckJSONTags checks whether field names match with JSON struct tags.
func CheckJSONTags(typ reflect.Type, terms ...string) error {
	return CheckStructTags("json", SnakeCase, typ, terms...)
}

// CheckXMLTags checks whether field names match with XML struct tags.
func CheckXMLTags(typ reflect.Type, terms ...string) error {
	return CheckStructTags("xml", SnakeCase, typ, terms...)
}
