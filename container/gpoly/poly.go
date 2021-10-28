package gpoly

import (
	"github.com/cryptowilliam/goutil/container/ggeometry"
	"github.com/cryptowilliam/goutil/container/gpoly/polyfit"
	"math"
)

/*
polyfit,polyval用法示例
x=[0.0 0.1 0.2 0.3 0.5 0.8 1.0]
y=[1.0 0.41 0.50 0.61 0.91 2.02 2.46]
A=polyfit(x,y,2);
z=polyval(A,x);
*/

// 多项式曲线拟合函数，也就是离散点拟合直线
// xs: 横坐标值
// ys：纵坐标值
// degree：拟合多项式的次数
// 输出：degree次拟合多项式系数(coefficients)，其长度=degree+1，第一个值就是斜率，后面的值是截距误差。若要计算多项式在横坐标上的值，还要继续调用函数PolyVal
func PolyFit(xs, ys []float64, degree int) []float64 {
	return polyfit.NewFitting(xs, ys, degree).Solve(true)
}

func PolyVal(xs []float64, coefficients []float64) []float64 {
	ys := make([]float64, len(xs))

	for i := 0; i < len(ys); i++ {
		ys[i] = evalPoly(xs[i], coefficients)
	}
	return ys
}

// x: 横坐标值
// coeffs: PolyFit计算出的系数
// Evaluate a polynomial at a point
func evalPoly(x float64, coeffs []float64) float64 {
	ret := 0.0

	for i, coefficient := range coeffs {
		ret += math.Pow(x, float64(i)) * coefficient
	}
	return ret
}

// 等间距离散点拟合直线的夹角
func PolyAngle(ys []float64) float64 {
	xs := make([]float64, len(ys))
	for i := 0; i < len(ys); i++ {
		xs[i] = float64(i + 1)
	}
	coffs := PolyFit(xs, ys, 2)
	return ggeometry.SlopeToAngle(coffs[0], coffs[1])
}
