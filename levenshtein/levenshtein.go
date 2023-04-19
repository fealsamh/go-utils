package levenshtein

// Matrix is a matrix.
type Matrix[T any] struct {
	data   []T
	height int
	width  int
}

// NewMatrix creates a matrix.
func NewMatrix[T any](h, w int) *Matrix[T] {
	return &Matrix[T]{
		data:   make([]T, h*w),
		height: h,
		width:  w,
	}
}

func (mx *Matrix[T]) Get(n, m int) T {
	return mx.data[n*mx.width+m]
}

func (mx *Matrix[T]) Set(n, m int, v T) {
	mx.data[n*mx.width+m] = v
}

// Action is the action taken at a position in the dynamic matrix.
type Action uint8

const (
	none Action = iota
	del
	ins
)

func (a Action) String() string {
	switch a {
	case 0:
		return "-"
	case 1:
		return "D"
	case 2:
		return "I"
	}
	return ""
}

type cell[T any] struct {
	dist int
	act  Action
}

// Distance calculates the Levenshtein distance between two sequences.
func Distance[T comparable](s, t []T) (int, []Action) {
	d := NewMatrix[cell[T]](len(s)+1, len(t)+1)
	for i := 1; i <= len(s); i++ {
		d.Set(i, 0, cell[T]{dist: i})
	}
	for j := 1; j <= len(t); j++ {
		d.Set(0, j, cell[T]{dist: j})
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			var substcost int
			if s[i-1] != t[j-1] {
				substcost = 1
			}
			idx, min := min(d.Get(i-1, j).dist+1, d.Get(i, j-1).dist+1, d.Get(i-1, j-1).dist+substcost)
			var act Action
			switch idx {
			case 0:
				act = del
			case 1:
				act = ins
			}
			d.Set(i, j, cell[T]{dist: min, act: act})
		}
	}
	i, j := len(s), len(t)
	var acts []Action
	for i > 0 || j > 0 {
		act := d.Get(i, j).act
		acts = append(acts, act)
		switch act {
		case del:
			i--
		case ins:
			j--
		case none:
			i--
			j--
		}
	}
	for i, j := 0, len(acts)-1; i < j; i, j = i+1, j-1 {
		acts[i], acts[j] = acts[j], acts[i]
	}
	return d.Get(len(s), len(t)).dist, acts
}

// Profiles calculates the distance between and differential profiles of two sequences.
func Profiles[T comparable](s, t []T, blank T) (int, []T, []T) {
	d, acts := Distance(s, t)
	var prof1, prof2 []T
	var i, j, k int
	for k < len(acts) {
		switch acts[k] {
		case none:
			prof1 = append(prof1, s[i])
			prof2 = append(prof2, t[j])
			i++
			j++
		case del:
			prof1 = append(prof1, s[i])
			prof2 = append(prof2, blank)
			i++
		case ins:
			prof1 = append(prof1, blank)
			prof2 = append(prof2, t[j])
			j++
		}
		k++
	}
	return d, prof1, prof2
}

func min(ns ...int) (int, int) {
	var idx int
	min := int(^uint(0) >> 1)
	for i, n := range ns {
		if n < min {
			min = n
			idx = i
		}
	}
	return idx, min
}
