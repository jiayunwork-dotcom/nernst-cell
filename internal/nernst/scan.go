package nernst

import (
	"nernst-cell/internal/phys"
)

// ScanPoint 是扫描序列中的一点。
type ScanPoint struct {
	Label               string
	EquilibriumPotentialV float64
	SlopeMVPerDecade    float64
}

// TemperatureScan 在给定摄氏温度序列上重算平衡电位与 Nernst
// 斜率。用于观察「只升温时斜率 RT/nF 变陡」的交叉规则：
// 返回点按输入温度顺序排列。
func TemperatureScan(in Input, temps []float64) ([]ScanPoint, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	out := make([]ScanPoint, 0, len(temps))
	for _, tc := range temps {
		cur := in
		cur.TemperatureC = tc
		e, err := EquilibriumPotential(cur)
		if err != nil {
			return nil, err
		}
		slope, err := SlopeMVPerDecade(cur)
		if err != nil {
			return nil, err
		}
		out = append(out, ScanPoint{
			Label:                 phys.FormatMV(tc) + " °C",
			EquilibriumPotentialV: e,
			SlopeMVPerDecade:      slope,
		})
	}
	return out, nil
}

// RatioScan 在给定活度比值序列上重算平衡电位。
// decades 以 decade 表示：值 0 表示 a_ox/a_red = 1，
// 值 1 表示比值 10，值 -1 表示比值 0.1。
func RatioScan(in Input, decades []float64) ([]ScanPoint, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	out := make([]ScanPoint, 0, len(decades))
	for _, d := range decades {
		ratio := phys.RatioAtDecades(d)
		e, err := EquilibriumAtRatio(in, ratio)
		if err != nil {
			return nil, err
		}
		slope, err := SlopeMVPerDecade(in)
		if err != nil {
			return nil, err
		}
		out = append(out, ScanPoint{
			Label:                 "10^" + phys.FormatRatio(d),
			EquilibriumPotentialV: e,
			SlopeMVPerDecade:      slope,
		})
	}
	return out, nil
}

// ScanTable 把扫描点格式化为文本表格。
func ScanTable(points []ScanPoint) string {
	var b []byte
	for _, p := range points {
		b = append(b, phys.PadLabel(p.Label, 12)...)
		b = append(b, ' ', '|', ' ')
		b = append(b, phys.FormatV(p.EquilibriumPotentialV)...)
		b = append(b, ' ', 'V', ' ', '|', ' ')
		b = append(b, phys.FormatMV(p.SlopeMVPerDecade)...)
		b = append(b, ' ', 'm', 'V', '/', 'd', 'e', 'c', '\n')
	}
	return string(b)
}
