package sexpr

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/fealsamh/go-utils/tokeniser"
)

// Parse parses an s-expression from a slice of tokens.
func Parse(tokens []tokeniser.Token) ([]interface{}, error) {
	list, _, err := parseList(tokens)
	return list, err
}

func parseList(tokens []tokeniser.Token) ([]interface{}, []tokeniser.Token, error) {
	t := tokens[0]
	if !(t.Type == tokeniser.Symbol && t.Value == "(") {
		return nil, nil, fmt.Errorf("expected '(' (line %d)", t.Location.Line)
	}
	tokens = tokens[1:]
	var list []interface{}
	for {
		t = tokens[0]
		if t.Type == tokeniser.Symbol && t.Value == ")" {
			return list, tokens[1:], nil
		}
		var (
			expr interface{}
			err  error
		)
		expr, tokens, err = parseElem(tokens)
		if err != nil {
			return nil, nil, err
		}
		list = append(list, expr)
	}
}

func parseElem(tokens []tokeniser.Token) (interface{}, []tokeniser.Token, error) {
	t := tokens[0]
	if t.Type == tokeniser.Symbol && t.Value == "(" {
		return parseList(tokens)
	}
	switch t.Type {
	case tokeniser.Ident:
		return t.Value, tokens[1:], nil
	case tokeniser.Int:
		n, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, nil, err
		}
		return n, tokens[1:], nil
	case tokeniser.Float:
		x, err := strconv.ParseFloat(t.Value, 64)
		if err != nil {
			return nil, nil, err
		}
		return x, tokens[1:], nil
	case tokeniser.String:
		return String(t.Value), tokens[1:], nil
	default:
		return nil, nil, fmt.Errorf("expected identifier, number, string or '(' (line %d)", t.Location.Line)
	}
}

// String is a string literal occurring in an s-expression.
type String string

// ParseString parses an s-expression.
func ParseString(code string) ([]interface{}, error) {
	tokens, err := tokeniser.Tokenise(bufio.NewReader(strings.NewReader(code)), &tokeniser.Config{
		IdentRunes:  []rune{'_', '-', '+', '/', '*', '<', '=', '>', '?', '!'},
		CommentRune: '#',
	})
	if err != nil {
		return nil, err
	}
	return Parse(tokens)
}

// Must panics if `err` isn't `nil`.
func Must[T any](data T, err error) T {
	if err != nil {
		panic(err)
	}
	return data
}
