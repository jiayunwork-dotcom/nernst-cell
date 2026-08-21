package nernst

import (
	"errors"
)

// errRatioNotPositive 在活度比非正时返回。
var errRatioNotPositive = errors.New("activity ratio must be positive")

// CellPolarity 描述一个浓差电池的核算结果。
type CellPolarity struct {
	// PotentialV 是电池电位，单位 V。
	PotentialV float64
	// ConcSidePositive 指示浓侧是否为正极。
	ConcSidePositive bool
	// ConcActivity 与 DilActivity 回显两侧活度。
	ConcActivity float64
	DilActivity  float64
}

// ConcentrationCell 计算浓差电池电位。
// 两极同种、E° 相消，电位只剩浓度比一项：
//
//	E = (RT/(nF))·ln(a_conc / a_dil)
//
// a_conc > a_dil 时 E 为正，浓侧为正极。Electrons、tempC 校验
// 同 Nernst 输入。
func ConcentrationCell(concActivity, dilActivity float64, electrons int, tempC float64) (CellPolarity, error) {
	in := Input{
		StandardPotentialV: 0,
		Electrons:          electrons,
		TemperatureC:       tempC,
		OxActivity:         concActivity,
		RedActivity:        dilActivity,
	}
	e, err := EquilibriumPotential(in)
	if err != nil {
		return CellPolarity{}, err
	}
	return CellPolarity{
		PotentialV:       e,
		ConcSidePositive: e > 0,
		ConcActivity:     concActivity,
		DilActivity:      dilActivity,
	}, nil
}
