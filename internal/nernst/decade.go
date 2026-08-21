package nernst

import (
	"math"

	"nernst-cell/internal/phys"
)

// PotentialShiftPerDecade 返回活度比 a_ox/a_red 乘 10 时
// 平衡电位的移动量，单位 V。正值表示比值增大电位升高。
// 等价于 SlopeVPerDecade。
func PotentialShiftPerDecade(in Input) (float64, error) {
	return SlopeVPerDecade(in)
}

// TenfoldRatioInputs 返回一对保持输入其余字段不变、
// 仅把活度比放大/缩小十倍的输入，便于做对照：
// ratio=10 表示 a_ox/a_red = 原比值×10。
func TenfoldRatioInputs(in Input, decades int) (Input, Input, error) {
	if err := Validate(in); err != nil {
		return Input{}, Input{}, err
	}
	scale := phys.RatioAtDecades(float64(decades))
	up := in
	up.OxActivity = in.OxActivity * scale
	down := in
	down.RedActivity = in.RedActivity * scale
	return up, down, nil
}

// EquilibriumDeltaDecades 返回给定输入从 ratio₁ 变到 ratio₂ 时
// 平衡电位的移动量（V），公式为 (RT/(nF))·ln(ratio₂/ratio₁)。
func EquilibriumDeltaDecades(in Input, ratio1, ratio2 float64) (float64, error) {
	if err := Validate(in); err != nil {
		return 0, err
	}
	if !phys.PositiveFinite(ratio1) || !phys.PositiveFinite(ratio2) {
		return 0, errRatioNotPositive
	}
	rtf := phys.ThermalVoltage(in.TemperatureC)
	return rtf / float64(in.Electrons) * (math.Log(ratio2) - math.Log(ratio1)), nil
}
