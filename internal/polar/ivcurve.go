package polar

import (
	"math"

	"nernst-cell/internal/nernst"
)

// EtaGrid 生成过电位网格 [min, max] 上均匀分布的 count 个点。
func EtaGrid(etaMinV, etaMaxV float64, count int) []float64 {
	grid := make([]float64, 0, count)
	if count <= 1 {
		return []float64{etaMinV}
	}
	for i := 0; i < count; i++ {
		f := float64(i) / float64(count-1)
		grid = append(grid, etaMinV+(etaMaxV-etaMinV)*f)
	}
	return grid
}

// BuildIVCurve 构造完整的极化曲线：由 Nernst 输入算平衡电位，
// 在过电位网格上逐点计算 BV 电流与 Tafel 参考电流。
func BuildIVCurve(in IVInput) (IVResult, error) {
	if err := ValidateIVInput(in); err != nil {
		return IVResult{}, err
	}
	eq, err := nernst.EquilibriumPotential(in.Nernst)
	if err != nil {
		return IVResult{}, err
	}
	slopeMV, err := nernst.SlopeMVPerDecade(in.Nernst)
	if err != nil {
		return IVResult{}, err
	}
	tafelMV, err := TafelSlopeMVPerDecade(in.AlphaAnodic, in.Nernst.TemperatureC)
	if err != nil {
		return IVResult{}, err
	}
	grid := EtaGrid(in.EtaMinV, in.EtaMaxV, in.EtaPoints)
	points := make([]IVPoint, 0, len(grid))
	for _, eta := range grid {
		ibv, err := Current(in.ExchangeCurrentDensity, in.AlphaAnodic, in.Nernst.Electrons, in.Nernst.TemperatureC, eta)
		if err != nil {
			return IVResult{}, err
		}
		itaf, err := TafelReference(in.ExchangeCurrentDensity, in.AlphaAnodic, in.Nernst.Electrons, in.Nernst.TemperatureC, eta)
		if err != nil {
			return IVResult{}, err
		}
		points = append(points, IVPoint{EtaV: eta, IBV: ibv, ITafel: itaf})
	}
	return IVResult{
		EquilibriumPotentialV:  eq,
		SlopeMVPerDecade:       slopeMV,
		TafelSlopeMVPerDecade:  tafelMV,
		ExchangeCurrentDensity: in.ExchangeCurrentDensity,
		AlphaAnodic:            in.AlphaAnodic,
		AlphaCathodic:          CathodicTransfer(in.AlphaAnodic, in.Nernst.Electrons),
		Points:                 points,
	}, nil
}

// MaxAbsCurrent 返回点列中 |i_BV| 的最大值，用于图表刻度。
func (r IVResult) MaxAbsCurrent() float64 {
	max := 0.0
	for _, p := range r.Points {
		if a := math.Abs(p.IBV); a > max {
			max = a
		}
	}
	return max
}

// AtEta 返回离给定过电位最近的点的 BV 电流。
func (r IVResult) AtEta(eta float64) (IVPoint, bool) {
	if len(r.Points) == 0 {
		return IVPoint{}, false
	}
	best := r.Points[0]
	bestDelta := math.Abs(r.Points[0].EtaV - eta)
	for _, p := range r.Points[1:] {
		if d := math.Abs(p.EtaV - eta); d < bestDelta {
			best = p
			bestDelta = d
		}
	}
	return best, true
}
