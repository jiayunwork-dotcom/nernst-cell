// Package polar 实现电极极化核算：由平衡电位出发定义过电位
// η = E − E_eq，按 Butler–Volmer 计算过电位下电流，并提供
// Tafel 斜率与 Tafel 极限参考电流用于对照。电流符号全程统一
// （阳极 η>0 电流为正），不在阴极/阳极之间对调两套约定。
package polar

import "nernst-cell/internal/nernst"

// IVInput 是一次极化曲线核算的输入。
// Nernst 字段给出平衡电位来源；ExchangeCurrentDensity 是交换
// 电流密度 i0（A/cm²，必须为正）；AlphaAnodic 是阳极传递系数
// α_a，阴极系数取 α_c = n − α_a。EtaMinV/EtaMaxV 定义过电位
// 网格范围，EtaPoints 是取点数。
type IVInput struct {
	Nernst                nernst.Input
	ExchangeCurrentDensity float64
	AlphaAnodic            float64
	EtaMinV                float64
	EtaMaxV                float64
	EtaPoints              int
}

// IVPoint 是极化曲线上一点。
type IVPoint struct {
	// EtaV 是过电位，单位 V。
	EtaV float64
	// IBV 是 Butler–Volmer 电流密度，单位 A/cm²。
	IBV float64
	// ITafel 是 Tafel 极限参考电流密度：η>0 时取
	// i0·exp(α_a F η/(RT))，η<0 时取 −i0·exp(α_c F |η|/(RT))。
	ITafel float64
}

// IVResult 是极化曲线核算结果。
type IVResult struct {
	// EquilibriumPotentialV 是由 Nernst 算出的平衡电位。
	EquilibriumPotentialV float64
	// SlopeMVPerDecade 是 Nernst 斜率，单位 mV/decade。
	SlopeMVPerDecade float64
	// TafelSlopeMVPerDecade 是 Tafel 斜率 b = 2.303RT/(αF)，mV/decade。
	TafelSlopeMVPerDecade float64
	// ExchangeCurrentDensity 回显交换电流密度。
	ExchangeCurrentDensity float64
	// AlphaAnodic 与 AlphaCathodic 回显传递系数。
	AlphaAnodic   float64
	AlphaCathodic float64
	// Points 是过电位网格上的点列。
	Points []IVPoint
}
