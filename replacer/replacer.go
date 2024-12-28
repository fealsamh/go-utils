package replacer

import (
	"strings"
)

// Replace replaces segments in the template according to the provided map.
func Replace(tmpl string, repl map[string]string) string {
	var (
		state int
		start int
		sb    strings.Builder
	)
	for i, c := range tmpl {
		switch state {
		case 0:
			if c == '{' {
				state = 1
				start = i
			} else {
				sb.WriteRune(c)
			}
		case 1:
			if c == '}' {
				seg := string(tmpl[start+1 : i])
				if t, ok := repl[seg]; ok {
					sb.WriteString(t)
				} else {
					sb.WriteString(seg)
				}
				state = 0
			}
		}
		i++
	}
	return sb.String()
}
