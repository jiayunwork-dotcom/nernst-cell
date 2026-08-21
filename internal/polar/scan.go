package polar

import (
	"nernst-cell/internal/nernst"
)

// CurrentScan 在显式过电位列表上逐点计算 BV 电流与 Tafel 参考
// 电流，返回与输入列表等长的点列。与 BuildIVCurve 的区别在于
// 网格由调用方指定，不要求等距。
func CurrentScan(in IVInput, etas []float64) (IVResult, error) {
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
	points := make([]IVPoint, 0, len(etas))
	for _, eta := range etas {
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

// AnodicBranch 取点列中 η > 0 的子集，用于单独观察阳极 Tafel 区。
func (r IVResult) AnodicBranch() []IVPoint {
	out := make([]IVPoint, 0, len(r.Points))
	for _, p := range r.Points {
		if p.EtaV > 0 {
			out = append(out, p)
		}
	}
	return out
}

// CathodicBranch 取点列中 η < 0 的子集。
func (r IVResult) CathodicBranch() []IVPoint {
	out := make([]IVPoint, 0, len(r.Points))
	for _, p := range r.Points {
		if p.EtaV < 0 {
			out = append(out, p)
		}
	}
	return out
}

// Point 返回点列第 i 个点；越界返回零值与 false。
func (r IVResult) Point(i int) (IVPoint, bool) {
	if i < 0 || i >= len(r.Points) {
		return IVPoint{}, false
	}
	return r.Points[i], true
}
