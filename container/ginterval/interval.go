package ginterval

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/container/gstring"
	"strings"
)

type (
	Interval struct {
		Min        *gdecimal.Decimal
		IncludeMin bool
		Max        *gdecimal.Decimal
		IncludeMax bool
	}
)

// Parse interval string: [1,2] [100,+∞) (-∞,-1)
func Parse(s string) (*Interval, error) {
	defErr := gerrors.New("invalid interval string %s", s)
	s = strings.Replace(s, " ", "", -1)
	ss := strings.Split(s, ",")
	if len(ss) != 2 {
		return nil, defErr
	}
	if !gstring.StartWith(ss[0], "[") && !gstring.StartWith(ss[0], "(") {
		return nil, defErr
	}
	if !gstring.EndWith(ss[1], "]") && !gstring.EndWith(ss[1], ")") {
		return nil, defErr
	}
	if strings.Contains(ss[0], "∞") && ss[0] != "(-∞" {
		return nil, defErr
	}
	if strings.Contains(ss[1], "∞") && ss[1] != "+∞)" {
		return nil, defErr
	}

	res := &Interval{}

	leftNum := gstring.RemoveHead(ss[0], 1)
	if leftNum == "-∞" {
		res.Min = nil
	} else if gnum.IsDigit(leftNum) {
		min, err := gdecimal.NewFromString(leftNum)
		if err != nil {
			return nil, defErr
		}
		res.Min = &min
	} else {
		return nil, defErr
	}

	rightNum := gstring.RemoveTail(ss[1], 1)
	if rightNum == "+∞" {
		res.Max = nil
	} else if gnum.IsDigit(rightNum) {
		max, err := gdecimal.NewFromString(rightNum)
		if err != nil {
			return nil, defErr
		}
		res.Max = &max
	} else {
		return nil, defErr
	}

	if res.Min != nil && res.Max != nil && res.Min.GreaterThan(*res.Max) {
		return nil, defErr
	}

	res.IncludeMin = res.Min != nil && gstring.StartWith(s, "[")
	res.IncludeMax = res.Max != nil && gstring.EndWith(s, "]")

	return res, nil
}

// New infinite interval (-∞,+∞)
func New() *Interval {
	return &Interval{}
}

// Set interval left edge.
func (i *Interval) SetMin(min gdecimal.Decimal, includeMin bool) *Interval {
	i.Min = &min
	i.IncludeMin = includeMin
	if i.Max != nil && i.Max.LessThan(*i.Min) {
		panic(gerrors.New("min %s > max %s", i.Min.String(), i.Max.String()))
	}
	return i
}

// Set interval right edge.
func (i *Interval) SetMax(max gdecimal.Decimal, includeMax bool) *Interval {
	i.Max = &max
	i.IncludeMax = includeMax
	if i.Min != nil && i.Min.GreaterThan(*i.Max) {
		panic(gerrors.New("min %s > max %s", i.Min.String(), i.Max.String()))
	}
	return i
}

// Check contains n or not.
func (i *Interval) Contains(n gdecimal.Decimal) bool {
	minCheck := true
	if i.Min != nil {
		if i.IncludeMin {
			minCheck = n.GreaterThanOrEqual(*i.Min)
		} else {
			minCheck = n.GreaterThan(*i.Min)
		}
	}

	maxCheck := true
	if i.Max != nil {
		if i.IncludeMax {
			maxCheck = n.LessThanOrEqual(*i.Max)
		} else {
			maxCheck = n.LessThan(*i.Max)
		}
	}

	return minCheck && maxCheck
}

// Check whether two intervals has overlapping area.
func (i *Interval) IsOverlap(cmp Interval) bool {
	// 'i' is on the left of 'cmp' entirely.
	if i.Max != nil && cmp.Min != nil {
		if i.IncludeMax && cmp.IncludeMin {
			if i.Max.LessThan(*cmp.Min) {
				return false
			}
		} else {
			if i.Max.LessThanOrEqual(*cmp.Min) {
				return false
			}
		}
	}

	// 'i' is on the right of 'cmp' entirely.
	if i.Min != nil && cmp.Max != nil {
		if i.IncludeMin && cmp.IncludeMax {
			if i.Min.GreaterThan(*cmp.Max) {
				return false
			}
		} else {
			if i.Min.GreaterThanOrEqual(*cmp.Max) {
				return false
			}
		}
	}

	return true
}

// Format interval to string.
func (i *Interval) String() string {
	s := ""

	if i.Min == nil {
		s = "(-∞"
	} else {
		if i.IncludeMin {
			s = "[" + i.Min.String()
		} else {
			s = "(" + i.Min.String()
		}
	}

	s += ","

	if i.Max == nil {
		s += "+∞)"
	} else {
		if i.IncludeMax {
			s += i.Max.String() + "]"
		} else {
			s += i.Max.String() + ")"
		}
	}

	return s
}
