package nernst

import (
	"math"

	"nernst-cell/internal/phys"
)

// Evaluate 执行一次完整的 Nernst 核算：校验输入、计算平衡
// 电位与 Nernst 斜率，并打包成 Result。
func Evaluate(in Input) (Result, error) {
	if err := Validate(in); err != nil {
		return Result{}, err
	}
	e, err := EquilibriumPotential(in)
	if err != nil {
		return Result{}, err
	}
	slopeMV, err := SlopeMVPerDecade(in)
	if err != nil {
		return Result{}, err
	}
	rtf := phys.ThermalVoltage(in.TemperatureC)
	decadeV := rtf / float64(in.Electrons) * math.Ln10
	return Result{
		EquilibriumPotentialV: e,
		SlopeMVPerDecade:      slopeMV,
		ThermalVoltageV:       rtf,
		DecadeShiftV:          decadeV,
		ActivityRatio:         in.OxActivity / in.RedActivity,
		TemperatureC:          in.TemperatureC,
	}, nil
}
