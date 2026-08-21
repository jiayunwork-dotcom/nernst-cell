package web

import "math"

// roundTo 把浮点数舍入到指定小数位，抑制 JSON 输出中的
// 浮点噪声（如 0.0591593… 显示为 0.059159）。
func roundTo(x float64, places int) float64 {
	scale := math.Pow(10, float64(places))
	return math.Round(x*scale) / scale
}

// roundE 把电位舍入到 6 位小数（微伏级展示）。
func roundE(v float64) float64 {
	return roundTo(v, 6)
}

// roundSlope 把斜率舍入到 3 位小数。
func roundSlope(v float64) float64 {
	return roundTo(v, 3)
}

// roundCurrent 把电流舍入到 8 位有效量级：绝对值小时保留
// 更多小数位，避免点列数值被抹平。
func roundCurrent(i float64) float64 {
	if i == 0 {
		return 0
	}
	return math.Round(i*1e6) / 1e6
}
