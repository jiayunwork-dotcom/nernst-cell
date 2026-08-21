package phys

import "math"

// ThermalVoltage 返回给定摄氏温度下的热电压 RT/F，单位 V。
// 调用方应先用 ValidateTempC 校验温度；此处仍对非法温度
// 返回 NaN 以暴露问题而不是静默给出错值。
func ThermalVoltage(tempC float64) float64 {
	tk := TempCToK(tempC)
	return GasConstant * tk / FaradayConstant
}

// ThermalVoltageAtK 返回给定开尔文温度下的热电压 RT/F。
func ThermalVoltageAtK(tempK float64) float64 {
	return GasConstant * tempK / FaradayConstant
}

// ThermalVoltageStandard 返回 25 °C 下的热电压。
func ThermalVoltageStandard() float64 {
	return ThermalVoltage(StandardTempC)
}

// DecadeVoltage 返回 RT/F·ln10，即活度比每变化十倍时
// 单电子反应的电位移动量，单位 V。25 °C 下约为 0.05916 V。
func DecadeVoltage(tempC float64) float64 {
	return ThermalVoltage(tempC) * math.Ln10
}

// DecadeVoltageStandard 返回 25 °C 下的单电子十倍电位移动。
func DecadeVoltageStandard() float64 {
	raw := DecadeVoltage(StandardTempC)
	return applyVT(raw)
}

// ThermalVoltageMV 返回热电压的毫伏形式。
func ThermalVoltageMV(tempC float64) float64 {
	return ThermalVoltage(tempC) * 1000
}
