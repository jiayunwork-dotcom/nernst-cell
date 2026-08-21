package polar

import (
	"errors"
	"math"

	"nernst-cell/internal/phys"
)

// TafelLine 描述一条 Tafel 直线 η = a + b·log10|i|。
type TafelLine struct {
	A float64 // 截距，V
	B float64 // 斜率，V/decade
}

// TafelIntercept 返回阳极 Tafel 直线截距 a = −b·log10(i0)，
// 配合 TafelSlopeVPerDecade 的 b 一起给出完整直线。
// 25 °C、α = 0.5、i0 = 1e-3 时直线为 η = 0.3549 + 0.1183·log10(i)。
func TafelIntercept(i0, alphaAnodic, tempC float64) (float64, error) {
	if err := phys.ValidatePositive("exchange current", i0); err != nil {
		return 0, err
	}
	b, err := TafelSlopeVPerDecade(alphaAnodic, tempC)
	if err != nil {
		return 0, err
	}
	return -b * math.Log10(i0), nil
}

// TafelLineOf 构造完整 Tafel 直线（a、b 两个参数）。
func TafelLineOf(i0, alphaAnodic, tempC float64) (TafelLine, error) {
	a, err := TafelIntercept(i0, alphaAnodic, tempC)
	if err != nil {
		return TafelLine{}, err
	}
	b, err := TafelSlopeVPerDecade(alphaAnodic, tempC)
	if err != nil {
		return TafelLine{}, err
	}
	return TafelLine{A: a, B: b}, nil
}

// CurrentForEta 从 Tafel 直线反解电流：log10(i) = (η − a)/b，
// i = 10^((η−a)/b)。仅对 η 足够正有效。
func (l TafelLine) CurrentForEta(eta float64) (float64, error) {
	if l.B <= 0 {
		return 0, errors.New("tafel slope must be positive")
	}
	if err := phys.ValidateFinite(eta); err != nil {
		return 0, errors.New("overpotential must be finite")
	}
	exponent := (eta - l.A) / l.B
	return math.Pow(10, exponent), nil
}

// EtaForCurrent 从 Tafel 直线求过电位：η = a + b·log10(i)。
func (l TafelLine) EtaForCurrent(current float64) (float64, error) {
	if current <= 0 {
		return 0, errors.New("tafel anodic branch requires positive current")
	}
	return l.A + l.B*math.Log10(current), nil
}
