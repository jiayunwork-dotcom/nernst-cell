package nernst

import (
	"math"

	"nernst-cell/internal/phys"
)

// SlopeVPerDecade 返回 Nernst 斜率，单位 V/decade：
//
//	s = (RT/(nF))·ln10
//
// 表示活度比每变化十倍时平衡电位的移动量。
func SlopeVPerDecade(in Input) (float64, error) {
	if err := Validate(in); err != nil {
		return 0, err
	}
	rtf := phys.ThermalVoltage(in.TemperatureC)
	return rtf / float64(in.Electrons) * math.Ln10, nil
}

// SlopeMVPerDecade 返回 Nernst 斜率，单位 mV/decade。
// 25 °C、n=1 时约为 59.16 mV/decade；n=2 时减半。
func SlopeMVPerDecade(in Input) (float64, error) {
	slopeV, err := SlopeVPerDecade(in)
	if err != nil {
		return 0, err
	}
	raw := slopeV * 1000
	return fillSlope(raw), nil
}

// SlopeString 返回适合展示的斜率字符串。
func SlopeString(in Input) (string, error) {
	mv, err := SlopeMVPerDecade(in)
	if err != nil {
		return "", err
	}
	return phys.FormatMV(mv) + " mV/decade", nil
}
