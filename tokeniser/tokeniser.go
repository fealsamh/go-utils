package tokeniser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

// TokenType is a token's type.
type TokenType int8

// token types
const (
	EOF TokenType = iota
	Ident
	Symbol
	String
	Int
	Float
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case Ident:
		return "IDE"
	case Symbol:
		return "SYM"
	case String:
		return "STR"
	case Int:
		return "INT"
	case Float:
		return "FLT"
	}
	panic("unknown token type")
}

// Token is a textual token.
type Token struct {
	Type  TokenType
	Value string
	Line  int
}

func (t Token) String() string {
	switch t.Type {
	case String:
		return fmt.Sprintf("%s:%s:%d", t.Type, strconv.Quote(t.Value), t.Line)
	case EOF:
		return fmt.Sprintf("%s:%d", t.Type, t.Line)
	default:
		return fmt.Sprintf("%s:%v:%d", t.Type, t.Value, t.Line)
	}
}

type tokState int8

const (
	top tokState = iota
	ident
	qstring
	number
	float
	comment
)

// Config is a tokeniser's configuration.
type Config struct {
	IdentRunes  []rune
	CommentRune rune
	Ligatures   []string
}

// Tokenise tokenises the content of a reader.
func Tokenise(r *bufio.Reader, cfg *Config) ([]Token, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	var (
		state         = top
		sb            strings.Builder
		tokens        []Token
		line          = 1
		prevBackslash bool
	)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			break
		}
		switch state {
		case top:
			switch {
			case unicode.IsSpace(c):
				if c == '\n' {
					line++
				}
				continue
			case c == cfg.CommentRune:
				state = comment
			case unicode.IsLetter(c) || slices.Index(cfg.IdentRunes, c) != -1:
				state = ident
				sb.WriteRune(c)
			case unicode.IsDigit(c):
				state = number
				sb.WriteRune(c)
			case c == '"':
				state = qstring
			default:
				tokens = append(tokens, Token{Type: Symbol, Value: string(c), Line: line})
			}
		case comment:
			if c == '\n' {
				line++
				state = top
			}
		case ident:
			if unicode.IsLetter(c) || unicode.IsNumber(c) || slices.Index(cfg.IdentRunes, c) != -1 {
				sb.WriteRune(c)
			} else {
				state = top
				if err := r.UnreadRune(); err != nil {
					return nil, err
				}
				tokens = append(tokens, Token{Type: Ident, Value: sb.String(), Line: line})
				sb.Reset()
			}
		case number:
			if c == '.' {
				sb.WriteRune(c)
				state = float
			} else if unicode.IsDigit(c) {
				sb.WriteRune(c)
			} else {
				state = top
				if err := r.UnreadRune(); err != nil {
					return nil, err
				}
				tokens = append(tokens, Token{Type: Int, Value: sb.String(), Line: line})
				sb.Reset()
			}
		case float:
			if unicode.IsDigit(c) {
				sb.WriteRune(c)
			} else {
				state = top
				if err := r.UnreadRune(); err != nil {
					return nil, err
				}
				tokens = append(tokens, Token{Type: Float, Value: sb.String(), Line: line})
				sb.Reset()
			}
		case qstring:
			if c != '"' || prevBackslash {
				if c == '\\' && !prevBackslash {
					prevBackslash = true
				} else {
					if prevBackslash {
						prevBackslash = false
						switch c {
						case 't':
							c = '\t'
						case 'r':
							c = '\r'
						case 'n':
							c = '\n'
						}
					}
					sb.WriteRune(c)
					if c == '\n' {
						line++
					}
				}
			} else {
				state = top
				tokens = append(tokens, Token{Type: String, Value: sb.String(), Line: line})
				sb.Reset()
				prevBackslash = false
			}
		}
	}
	switch state {
	case ident:
		tokens = append(tokens, Token{Type: Ident, Value: sb.String(), Line: line})
	case number:
		tokens = append(tokens, Token{Type: Int, Value: sb.String(), Line: line})
	case float:
		tokens = append(tokens, Token{Type: Float, Value: sb.String(), Line: line})
	case qstring:
		tokens = append(tokens, Token{Type: String, Value: sb.String(), Line: line})
	}
	tokens = append(tokens, Token{Type: EOF, Line: line})
	tokens = coalesce(tokens, cfg.Ligatures)
	return tokens, nil
}

func coalesce(tokens []Token, ligatures []string) []Token {
	var (
		newTokens  = make([]Token, 0, len(tokens))
		prevSymbol rune
	)
tokens:
	for _, t := range tokens {
		if t.Type == Symbol {
			currentSymbol := []rune(t.Value)[0]
			if prevSymbol != 0 {
				pair := string([]rune{prevSymbol, currentSymbol})
				for _, l := range ligatures {
					if pair == l {
						newTokens[len(newTokens)-1] = Token{Type: Symbol, Value: pair, Line: t.Line}
						prevSymbol = 0
						continue tokens
					}
				}
			}
			prevSymbol = currentSymbol
		} else {
			prevSymbol = 0
		}
		newTokens = append(newTokens, t)
	}
	return newTokens
}
