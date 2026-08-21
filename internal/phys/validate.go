package phys

import (
	"errors"
	"math"
)

// PositiveFinite 判断 x 是否为正且有限。
// NaN 与 0 与负数返回 false；+Inf 返回 false。
func PositiveFinite(x float64) bool {
	return x > 0 && !math.IsInf(x, 0) && !math.IsNaN(x)
}

// NonNegativeFinite 判断 x 是否非负且有限。
func NonNegativeFinite(x float64) bool {
	return x >= 0 && !math.IsInf(x, 0) && !math.IsNaN(x)
}

// Finite 判断 x 是否有限（非 NaN 非 Inf）。
func Finite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}

// ErrMustBeFinite 是数值必须有限的通用错误。
var ErrMustBeFinite = errors.New("value must be finite")

// ValidateFinite 校验 x 有限，否则返回 ErrMustBeFinite。
func ValidateFinite(x float64) error {
	if !Finite(x) {
		return ErrMustBeFinite
	}
	return nil
}

// ValidatePositive 校验 x 为正且有限。
func ValidatePositive(name string, x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return errors.New(name + " must be finite")
	}
	if x <= 0 {
		return errors.New(name + " must be positive")
	}
	return nil
}

// ValidateElectronCount 校验电子数 n：必须为正整数。
// n 接受 float64 形式（0.5 不合法），便于从 JSON 直接校验。
func ValidateElectronCount(n float64) error {
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return errors.New("electrons must be finite")
	}
	if n <= 0 {
		return errors.New("electrons must be a positive integer")
	}
	if math.Abs(n-math.Round(n)) > 1e-9 {
		return errors.New("electrons must be a positive integer")
	}
	return nil
}
