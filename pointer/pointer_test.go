package pointer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {
	req := require.New(t)

	p1 := To("abcd")
	req.Equal("abcd", *p1)

	now := time.Now()
	p2 := To(now)
	req.Equal(now, *p2)

	id := uuid.New()
	p3 := To(id)
	req.Equal(id, *p3)
}
