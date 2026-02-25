package embedding

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"

	"github.com/fealsamh/go-utils/nocopy"
)

// Vector is a numerical vector with driver value support.
type Vector []float64

// Value ...
func (v Vector) Value() (driver.Value, error) {
	b := []byte{'['}
	for i, x := range v {
		if i > 0 {
			b = append(b, ',')
		}
		b = strconv.AppendFloat(b, x, 'E', -1, 64)
	}
	return append(b, ']'), nil
}

// Scan ...
func (v *Vector) Scan(src any) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("vector from database not slice of bytes")
	}
	if len(b) < 2 {
		return errors.New("vector from database invalid")
	}
	if b[0] != '[' || b[len(b)-1] != ']' {
		return errors.New("vector from database ill-formed")
	}
	var xs []float64
	for s := range strings.SplitSeq(nocopy.String(b[1:len(b)-1]), ",") {
		x, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		xs = append(xs, x)
	}
	*v = xs
	return nil
}
