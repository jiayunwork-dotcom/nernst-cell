package phys

import "math"

// Decades 返回以 10 为底的比值 decade 数：log10(a/b)。
// a 与 b 都应为正；非法输入返回 NaN。
func Decades(a, b float64) float64 {
	if !PositiveFinite(a) || !PositiveFinite(b) {
		return math.NaN()
	}
	return math.Log10(a / b)
}

// RatioAtDecades 返回 decade 数 d 对应的比值 10^d。
func RatioAtDecades(d float64) float64 {
	return math.Pow(10, d)
}

// DecadeFraction 把比值 a/b 表示成「多少倍十倍」。
// 例：a/b = 10 → 1.0；a/b = 1 → 0。
func DecadeFraction(a, b float64) float64 {
	return math.Log10(a / b)
}

// ScaleToDecade 返回把 x 调整到目标 decade 值所需的倍率。
// 用于构造「活度比乘十倍」的对照输入。
func ScaleToDecade(x float64, decades float64) float64 {
	return x * math.Pow(10, decades)
}
