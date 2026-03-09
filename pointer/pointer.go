package pointer

import (
	"time"

	"github.com/google/uuid"
)

type simpleType interface {
	int | uint | int32 | uint32 | int64 | uint64 | bool | float32 | float64 | string | time.Time | uuid.UUID
}

// To returns a pointer to its argument.
func To[T simpleType](x T) *T {
	return &x
}
