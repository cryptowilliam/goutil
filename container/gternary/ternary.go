package gternary

type CondExpr bool

func If(flag bool) CondExpr {
	return CondExpr(flag)
}

func (ce CondExpr) Int8(yes, no int8) int8 {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Int(yes, no int) int {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Int64(yes, no int64) int64 {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Uint8(yes, no uint8) uint8 {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Uint(yes, no uint) uint {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Uint64(yes, no uint64) uint64 {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Float64(yes, no float64) float64 {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) String(yes, no string) string {
	if ce {
		return yes
	}
	return no
}
func (ce CondExpr) Interface(yes, no interface{}) interface{} {
	if ce {
		return yes
	}
	return no
}
