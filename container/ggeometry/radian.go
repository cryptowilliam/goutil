package ggeometry

import "math"

// 斜率到弧度
// [0,0]和[x,y]这两个点的连线的弧度
func SlopeToRadian(x, y float64) float64 {
	return math.Atan2(x, y)
}

func RadianToAngle(x float64) float64 {
	return x * (180 / math.Pi)
}

func AngleToRadian(x float64) float64 {
	return x * (math.Pi / 180)
}

// TODO: test required
// [0,0]和[x,y]这两个点的连线的夹角
func SlopeToAngle(x, y float64) float64 {
	return RadianToAngle(SlopeToRadian(x, y))
}
