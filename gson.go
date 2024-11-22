package gson

/*
 * JSON データを Unmarshal したいが，
 * そのデータが null や 0 を指定されたのか undefined かを判定したいときに使う
 * int データは，どちらのケースでも 0 になり，区別ができない
 */

import (
	"encoding/json"
	"fmt"
)

type Gson[T any] struct {
	value   T
	defined bool
	null    bool
}

func (g *Gson[T]) UnmarshalJSON(data []byte) error {
	g.defined = true
	if string(data) == "null" {
		g.null = true
		return nil
	}
	return json.Unmarshal(data, &g.value)
}

func (g *Gson[T]) Value() T {
	return g.value
}

func (g *Gson[T]) IsDefined() bool {
	return g.defined
}

func (g *Gson[T]) IsNull() bool {
	return g.null
}

func (g *Gson[T]) Scan(src any) error {
	g.defined = true
	if src == nil {
		g.null = true
		return nil
	}
	switch v := src.(type) {
	case T:
		g.value = v
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return nil
}

func (g *Gson[T]) formatString(s fmt.State, verb rune) string {
	ff := "%"
	for _, f := range " -+#" {
		if s.Flag(int(f)) {
			ff += string(f)
		}
	}
	if w, ok := s.Width(); ok {
		ff += fmt.Sprintf("%d", w)
	}
	if p, ok := s.Precision(); ok {
		ff += fmt.Sprintf(".%d", p)
	}
	ff += string(verb)
	return ff
}

func (g Gson[T]) Format(s fmt.State, verb rune) {
	if g.null {
		str := "null"
		if w, ok := s.Width(); ok && w < 4 {
			if w <= 2 {
				str = "nl"
			} else {
				str = "nul"
			}
		}
		s.Write([]byte(str))
	} else if !g.defined {
		str := "---------"
		if w, ok := s.Width(); ok && w <= 5 {
			if w <= 2 {
				str = str[:w]
			} else if w <= 4 {
				str = str[:w]
			} else {
				str = str[:w]
			}
		}
		s.Write([]byte(str))
	} else {
		switch v := any(g.value).(type) {
		case bool, int, uint, string, []byte:
			s.Write([]byte(fmt.Sprintf(g.formatString(s, verb), v)))
		case fmt.Formatter:
			v.Format(s, verb)
		default:
			s.Write([]byte(fmt.Sprintf("%v", g.value)))
		}
	}
}
