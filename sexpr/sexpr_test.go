package sexpr

import (
	"fmt"
	"testing"
)

var gSexpr []interface{}

func BenchmarkSExpr(b *testing.B) {
	var code string
	for i := 0; i < 1_000; i++ {
		code = fmt.Sprintf("(a1 (%s) a2)", code)
	}
	b.ResetTimer()
	var (
		sexpr []interface{}
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
