package phys

// CurrentDensityUnits 是电流密度的展示单位。
type CurrentDensityUnits string

const (
	// AmpPerCm2 单位 A/cm²。
	AmpPerCm2 CurrentDensityUnits = "A/cm^2"
	// MilliAmpPerCm2 单位 mA/cm²。
	MilliAmpPerCm2 CurrentDensityUnits = "mA/cm^2"
)

// ToMilli 把 A/cm² 换算为 mA/cm²。
func ToMilli(i float64) float64 {
	return i * 1000
}

// FromMilli 把 mA/cm² 换算为 A/cm²。
func FromMilli(imA float64) float64 {
	return imA / 1000
}

// FormatCurrentWithUnit 按给定单位格式化电流密度。
func FormatCurrentWithUnit(i float64, unit CurrentDensityUnits) string {
	switch unit {
	case MilliAmpPerCm2:
		return FormatCurrent(ToMilli(i)) + " mA/cm^2"
	default:
		return FormatCurrent(i) + " A/cm^2"
	}
}

// ScaleCurrent 把电流密度缩放到易读量级：
// 绝对值 ≥ 1 用 A/cm²，否则用 mA/cm²。
func ScaleCurrent(i float64) (float64, CurrentDensityUnits) {
	if absf(i) >= 1 {
		return i, AmpPerCm2
	}
	return ToMilli(i), MilliAmpPerCm2
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
