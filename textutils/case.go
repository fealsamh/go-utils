package textutils

import (
	"strings"
	"unicode"
)

// SplitCamelcasedString splits up a camelcased string into components.
func SplitCamelcasedString(s string, terms ...string) []string {
	var (
		comps  []string
		start  int
		locked bool
	)
	runes := []rune(s)
	for i, r := range runes {
		switch {
		case i-start == 1 && unicode.IsUpper(r):
			locked = true
		case locked && unicode.IsLower(r):
			comps = append(comps, string(runes[start:i-1]))
			start = i - 1
			locked = false
		case i > 0 && !locked && unicode.IsUpper(r):
			comps = append(comps, string(runes[start:i]))
			start = i
		}
	}
	comps = append(comps, string(runes[start:]))
	r := comps
	for _, t := range terms {
		c := r
		r = make([]string, 0, len(c))
		for i := 0; i < len(c); i++ {
			if i < len(c)-1 {
				c1, c2 := strings.ToLower(c[i]), strings.ToLower(c[i+1])
				if c1 == t[:len(t)-1] && c2[:1] == t[len(t)-1:] {
					r = append(r, c1+c2)
					i++
					continue
				}
			}
			r = append(r, c[i])
		}
	}
	return r
}

// SnakecasedFromCamelcased converts a camelcased string into a snakecased string with lowercased components.
func SnakecasedFromCamelcased(s string, terms ...string) string {
	comps := SplitCamelcasedString(s, terms...)
	return strings.ToLower(strings.Join(comps, "_"))
}
