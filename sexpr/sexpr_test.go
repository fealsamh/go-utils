package sexpr

import (
	"fmt"
	"testing"
)

var gSexpr []any

func BenchmarkSExpr(b *testing.B) {
	var code string
	for range 1_000 {
		code = fmt.Sprintf("(a1 (%s) a2)", code)
	}
	b.ResetTimer()
	var (
		sexpr []any
		err   error
	)
	for n := 0; n < b.N; n++ {
		sexpr, err = ParseString(code)
		if err != nil {
			b.FailNow()
		}
	}
	b.Log(len(sexpr))
	gSexpr = sexpr
}
