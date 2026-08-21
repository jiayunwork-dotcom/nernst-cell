package phys

// Interval 表示闭区间 [Lo, Hi]，用于网格与范围的统一校验。
type Interval struct {
	Lo float64
	Hi float64
}

// NewInterval 构造区间；lo > hi 时交换两者，保证不变量。
func NewInterval(lo, hi float64) Interval {
	if lo > hi {
		lo, hi = hi, lo
	}
	return Interval{Lo: lo, Hi: hi}
}

// Contains 判断 x 是否在闭区间内。
func (iv Interval) Contains(x float64) bool {
	return x >= iv.Lo && x <= iv.Hi
}

// ContainsOpen 判断 x 是否在开区间 (lo, hi) 内。
func (iv Interval) ContainsOpen(x float64) bool {
	return x > iv.Lo && x < iv.Hi
}

// Span 返回区间长度 hi − lo。
func (iv Interval) Span() float64 {
	return iv.Hi - iv.Lo
}

// Mid 返回区间中点。
func (iv Interval) Mid() float64 {
	return (iv.Lo + iv.Hi) / 2
}

// Clamp 把 x 夹到区间内。
func (iv Interval) Clamp(x float64) float64 {
	if x < iv.Lo {
		return iv.Lo
	}
	if x > iv.Hi {
		return iv.Hi
	}
	return x
}

// InOpenUnitInterval 判断 x 是否在 (0, 1) 内。
func InOpenUnitInterval(x float64) bool {
	return x > 0 && x < 1
}

// InOpenInterval 判断 x 是否在 (lo, hi) 内。
func InOpenInterval(x, lo, hi float64) bool {
	return x > lo && x < hi
}
