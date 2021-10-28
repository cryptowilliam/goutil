package gstring

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
)

type (
	Checker struct {
		rules [][2]rune
	}
)

func NewChecker() *Checker {
	return &Checker{}
}

func (r *Checker) Allow(min, max rune) *Checker {
	if min > max {
		min, max = max, min
	}
	r.rules = append(r.rules, [2]rune{min, max})
	return r
}

func (r *Checker) AllowRune(a rune) *Checker {
	r.rules = append(r.rules, [2]rune{a, a})
	return r
}

func (r *Checker) String() string {
	var ss []string
	for _, rule := range r.rules {
		if rule[0] == rule[1] {
			ss = append(ss, fmt.Sprintf("[%s]", string(rule[0])))
		} else {
			ss = append(ss, fmt.Sprintf("[%s,%s]", string(rule[0]), string(rule[1])))
		}
	}
	return `{` + strings.Join(ss, " ") + `}`
}

func (r *Checker) Check(s string) error {
	for _, letter := range s {
		ok := false
		for _, rule := range r.rules {
			if letter >= rule[0] && letter <= rule[1] {
				ok = true
				break
			}
		}
		if !ok {
			return gerrors.New("string(%s) check error under rules: %s", s, r.String())
		}
	}
	return nil
}
