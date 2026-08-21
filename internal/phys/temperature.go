package phys

import (
	"errors"
	"math"
)

// KelvinZeroC 是 0 °C 对应的开尔文温度。
const KelvinZeroC = 273.15

// StandardTempC 是电化学习惯采用的 25 °C。
const StandardTempC = 25.0

// StandardTempK 是 25 °C 对应的开尔文温度。
const StandardTempK = KelvinZeroC + StandardTempC

// TempCToK 把摄氏温度换算为开尔文。
// 不校验输入；校验请走 ValidateTempC。
func TempCToK(tempC float64) float64 {
	return tempC + KelvinZeroC
}

// KToTempC 把开尔文温度换算为摄氏。
func KToTempC(tempK float64) float64 {
	return tempK - KelvinZeroC
}

// ValidateTempC 校验摄氏温度：
// 非有限值、T ≤ 0 都视为非法并返回带说明的 error。
// T 以摄氏度计，T = 0 表示冰点以下，对多数电化学过程无意义。
func ValidateTempC(tempC float64) error {
	if math.IsNaN(tempC) || math.IsInf(tempC, 0) {
		return errors.New("temperature must be finite")
	}
	if tempC <= 0 {
		return errors.New("temperature must be positive (celsius)")
	}
	return nil
}
