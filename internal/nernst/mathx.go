package nernst

import "math"

// abs 返回浮点绝对值，避免在包内多处重复 math.Abs。
func abs(x float64) float64 {
	return math.Abs(x)
}

// closeTo 判断两个浮点数是否在绝对容差内相等。
func closeTo(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}
