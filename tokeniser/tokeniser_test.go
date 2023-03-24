package tokeniser

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

var gTokens []*Token

func BenchmarkTokeniser(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 1_000; i++ {
		sb.WriteString(fmt.Sprintf("id%d ", i))
	}
	code := sb.String()
	b.ResetTimer()
	var (
		tokens []*Token
		err    error
	)
	for n := 0; n < b.N; n++ {
		tokens, err = Tokenise(bufio.NewReader(strings.NewReader(code)), nil)
		if err != nil {
			b.FailNow()
		}
	}
	b.Log(len(tokens))
	gTokens = tokens
}
