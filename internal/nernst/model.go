// Package nernst 实现 Nernst 平衡电位核算：给定标准电位、
// 电子数、温度与氧化/还原活度，计算平衡电位与 Nernst 斜率，
// 并支持浓差电池口径（E° 相消只剩浓度比）。
package nernst

// Input 是一次 Nernst 核算的输入。
// StandardPotentialV 是标准还原电位 E°（V，以还原反应式书写）；
// OxActivity 与 RedActivity 是氧化态/还原态活度（正数）；
// TemperatureC 以摄氏计，必须为正。
type Input struct {
	StandardPotentialV float64
	Electrons          int
	TemperatureC       float64
	OxActivity         float64
	RedActivity        float64
}

// Result 是一次 Nernst 核算的输出。
type Result struct {
	// EquilibriumPotentialV 是平衡电位 E，单位 V。
	EquilibriumPotentialV float64
	// SlopeMVPerDecade 是 Nernst 斜率，单位 mV/decade，
	// 即活度比每变十倍时电位移动量。
	SlopeMVPerDecade float64
	// ThermalVoltageV 是当前温度下的 RT/F。
	ThermalVoltageV float64
	// DecadeShiftV 是单电子十倍浓度比对应的电位移动，
	// 等于 SlopeMVPerDecade/1000。
	DecadeShiftV float64
	// ActivityRatio 是 a_ox / a_red。
	ActivityRatio float64
	// TemperatureC 回显输入温度。
	TemperatureC float64
}
