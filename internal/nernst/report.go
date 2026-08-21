package nernst

import (
	"nernst-cell/internal/phys"
)

// ResultRow 返回一行制表符分隔的结果，便于批量离线对比：
//
//	E  slope  RT/F  decade_shift  ratio  T
func (r Result) ResultRow() string {
	return phys.FormatV(r.EquilibriumPotentialV) + "\t" +
		phys.FormatMV(r.SlopeMVPerDecade) + "\t" +
		phys.FormatV(r.ThermalVoltageV) + "\t" +
		phys.FormatV(r.DecadeShiftV) + "\t" +
		phys.FormatRatio(r.ActivityRatio) + "\t" +
		phys.FormatMV(r.TemperatureC)
}

// ResultHeader 返回与 ResultRow 对齐的表头。
func ResultHeader() string {
	return "E(V)\t" +
		"slope(mV/dec)\t" +
		"RT/F(V)\t" +
		"decade_shift(V)\t" +
		"ratio\t" +
		"T(C)"
}

// Compare 比较两个 Result 的平衡电位差（V），正值表示前者更高。
func Compare(a, b Result) float64 {
	return a.EquilibriumPotentialV - b.EquilibriumPotentialV
}

// SameEquilibrium 判断两个 Result 的平衡电位在容差内相等。
func SameEquilibrium(a, b Result, tolV float64) bool {
	return closeTo(a.EquilibriumPotentialV, b.EquilibriumPotentialV, tolV)
}
