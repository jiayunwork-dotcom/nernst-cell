package nernst

import (
	"math"

	"nernst-cell/internal/phys"
)

// EquilibriumPotential 按 Nernst 方程计算平衡电位：
//
//	E = E° + (RT/(nF))·ln(a_ox / a_red)
//
// 还原态 a_red 写在分母，与标准还原电位表（还原反应式）一致。
// 输入非法时返回 error，不做静默错值。
func EquilibriumPotential(in Input) (float64, error) {
	if err := Validate(in); err != nil {
		return 0, err
	}
	rtf := phys.ThermalVoltage(in.TemperatureC)
	ratio := in.OxActivity / in.RedActivity
	term := rtf / float64(in.Electrons) * math.Log(ratio)
	return in.StandardPotentialV + term, nil
}

// EquilibriumAtRatio 返回把活度比设为指定 ratio 时的平衡电位。
// 用于「浓度比乘十倍」类对照计算：ratio 可以是 10、0.1 等。
func EquilibriumAtRatio(in Input, ratio float64) (float64, error) {
	if err := Validate(in); err != nil {
		return 0, err
	}
	if !phys.PositiveFinite(ratio) {
		return 0, errRatioNotPositive
	}
	rtf := phys.ThermalVoltage(in.TemperatureC)
	term := rtf / float64(in.Electrons) * math.Log(ratio)
	return in.StandardPotentialV + term, nil
}
