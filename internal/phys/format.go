package phys

import (
	"math"
	"strconv"
	"strings"
)

// FormatV 把伏特值格式化为适合展示的字符串：
// 绝对值小于 1e-4 用科学计数，否则保留 6 位有效数字。
func FormatV(v float64) string {
	if !Finite(v) {
		return "NaN"
	}
	if math.Abs(v) < 1e-4 && v != 0 {
		return strconv.FormatFloat(v, 'e', 2, 64)
	}
	return strconv.FormatFloat(v, 'f', 6, 64)
}

// FormatMV 把毫伏值格式化。
func FormatMV(mv float64) string {
	if !Finite(mv) {
		return "NaN"
	}
	return strconv.FormatFloat(mv, 'f', 2, 64)
}

// FormatCurrent 把电流密度格式化，单位 A/cm²：
// 绝对值过大或过小走科学计数。
func FormatCurrent(i float64) string {
	if !Finite(i) {
		return "NaN"
	}
	ai := math.Abs(i)
	if ai >= 1e4 || (ai < 1e-3 && ai != 0) {
		return strconv.FormatFloat(i, 'e', 3, 64)
	}
	return strconv.FormatFloat(i, 'f', 6, 64)
}

// FormatRatio 把无量纲比值格式化为 4 位有效数字。
func FormatRatio(r float64) string {
	if !Finite(r) {
		return "NaN"
	}
	return strconv.FormatFloat(r, 'g', 4, 64)
}

// FormatPercent 把 [0,1] 区间分数格式化为百分数。
func FormatPercent(frac float64) string {
	if !Finite(frac) {
		return "NaN"
	}
	return strconv.FormatFloat(frac*100, 'f', 2, 64) + "%"
}

// PadLabel 把标签补齐到固定宽度，用于表格对齐。
func PadLabel(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
