package nernst

import (
	"math"

	"nernst-cell/internal/phys"
)

// ConcentrationCellAtDecades 计算浓差电池在两侧活度相差给定
// decade 数时的电位。decades 为正表示浓侧活度是稀侧的 10^decades
// 倍，电池电位为正，浓侧为正极。decades = 0 时电位为 0。
func ConcentrationCellAtDecades(dilActivity float64, decades int, electrons int, tempC float64) (CellPolarity, error) {
	if err := phys.ValidatePositive("dilute activity", dilActivity); err != nil {
		return CellPolarity{}, err
	}
	conc := dilActivity * phys.RatioAtDecades(float64(decades))
	return ConcentrationCell(conc, dilActivity, electrons, tempC)
}

// CellLadder 返回两个相差 d_decades 个数量级的浓差电池电位，
// 用于演示每十倍浓度比移动 RT/nF·ln10。
func CellLadder(concActivity, dilActivity float64, electrons int, tempC float64, dDecades int) ([]ScanPoint, error) {
	if err := phys.ValidatePositive("concentrated activity", concActivity); err != nil {
		return nil, err
	}
	if err := phys.ValidatePositive("dilute activity", dilActivity); err != nil {
		return nil, err
	}
	baseRatio := concActivity / dilActivity
	out := make([]ScanPoint, 0, dDecades+1)
	for d := 0; d <= dDecades; d++ {
		ratio := baseRatio * math.Pow(10, float64(d))
		in := Input{
			StandardPotentialV: 0,
			Electrons:          electrons,
			TemperatureC:       tempC,
			OxActivity:         concActivity * math.Pow(10, float64(d)),
			RedActivity:        dilActivity,
		}
		e, err := EquilibriumPotential(in)
		if err != nil {
			return nil, err
		}
		slope, err := SlopeMVPerDecade(in)
		if err != nil {
			return nil, err
		}
		out = append(out, ScanPoint{
			Label:                 "ratio " + phys.FormatRatio(ratio),
			EquilibriumPotentialV: e,
			SlopeMVPerDecade:      slope,
		})
	}
	return out, nil
}
