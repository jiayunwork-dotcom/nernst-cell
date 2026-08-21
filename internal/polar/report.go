package polar

import (
	"strings"

	"nernst-cell/internal/phys"
)

// String 返回极化曲线结果的文本摘要，便于日志与离线查看。
func (r IVResult) String() string {
	var b strings.Builder
	b.WriteString("IV curve: E_eq = " + phys.FormatV(r.EquilibriumPotentialV) + " V")
	b.WriteString(", b = " + phys.FormatMV(r.TafelSlopeMVPerDecade) + " mV/decade")
	b.WriteString(", i0 = " + phys.FormatCurrent(r.ExchangeCurrentDensity) + " A/cm^2")
	b.WriteString(", alpha_a = " + phys.FormatRatio(r.AlphaAnodic))
	b.WriteString(", alpha_c = " + phys.FormatRatio(r.AlphaCathodic))
	if len(r.Points) > 0 {
		first := r.Points[0]
		last := r.Points[len(r.Points)-1]
		b.WriteString(", eta range [" + phys.FormatV(first.EtaV) + ", " + phys.FormatV(last.EtaV) + "] V")
		b.WriteString(", " + intString(len(r.Points)) + " points")
	}
	return b.String()
}

// PointTable 返回点列的等宽文本表格（制表符分隔），
// 每行形如：eta i_bv i_tafel。
func (r IVResult) PointTable() string {
	var b strings.Builder
	b.WriteString("eta_V\t i_bv\t i_tafel\n")
	for _, p := range r.Points {
		b.WriteString(phys.FormatV(p.EtaV))
		b.WriteString("\t")
		b.WriteString(phys.FormatCurrent(p.IBV))
		b.WriteString("\t")
		b.WriteString(phys.FormatCurrent(p.ITafel))
		b.WriteString("\n")
	}
	return b.String()
}

func intString(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return phys.FormatRatio(float64(n))
}
