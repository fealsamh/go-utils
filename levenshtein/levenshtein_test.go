package levenshtein

import "testing"

func TestDistance(t *testing.T) {
	cases := []struct {
		s  string
		t  string
		d  int
		p1 string
		p2 string
	}{
		{
			s:  "aaa",
			t:  "aaa",
			d:  0,
			p1: "aaa",
			p2: "aaa",
		},
		{
			s:  "aaa",
			t:  "aba",
			d:  1,
			p1: "aaa",
			p2: "aba",
		},
		{
			s:  "acd",
			t:  "abcd",
			d:  1,
			p1: "a cd",
			p2: "abcd",
		},
		{
			s:  "abcd",
			t:  "acd",
			d:  1,
			p1: "abcd",
			p2: "a cd",
		},
	}
	for _, c := range cases {
		d, prof1, prof2 := Profiles([]rune(c.s), []rune(c.t), ' ')
		if d != c.d {
			t.Errorf("expected distance between '%s' and '%s' is %d, got %d", c.s, c.t, c.d, d)
		}
		if string(prof1) != c.p1 {
			t.Errorf("expected first profile to be '%s', got %d", c.p1, prof1)
		}
		if string(prof2) != c.p2 {
			t.Errorf("expected second profile to be '%s', got %d", c.p2, prof2)
		}
	}
}
