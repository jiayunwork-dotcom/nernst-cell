package nernst

import (
	"nernst-cell/internal/phys"
)

// Summary 返回一次 Nernst 核算的多行文本摘要，
// 便于离线查看与日志输出。
func (r Result) Summary() string {
	return "E = " + phys.FormatV(r.EquilibriumPotentialV) + " V" +
		"\nNernst slope = " + phys.FormatMV(r.SlopeMVPerDecade) + " mV/decade" +
		"\nRT/F = " + phys.FormatV(r.ThermalVoltageV) + " V" +
		"\nsingle-electron decade shift = " + phys.FormatV(r.DecadeShiftV) + " V" +
		"\na_ox / a_red = " + phys.FormatRatio(r.ActivityRatio) +
		"\nT = " + phys.FormatMV(r.TemperatureC) + " °C"
}

// StandardShift 返回 25 °C、n=1 时十倍活度比对应的电位移动
// （V），即教科书上的 59.16 mV。作为交叉校验的常量入口。
func StandardShift() float64 {
	return phys.DecadeVoltageStandard()
}

// IsConcentrationCell 判断一次核算是否近似为浓差电池：
// 标准电位为零（或小到忽略）即视为浓度比主导。
func (r Result) IsConcentrationCell(tolV float64) bool {
	return abs(r.EquilibriumPotentialV) > tolV && r.ThermalVoltageV > 0
}
