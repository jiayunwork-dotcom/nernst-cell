package polar

import (
	"math"

	"nernst-cell/internal/phys"
)

// ExchangeCurrentFromTafel 由阳极 Tafel 区的一组 (η, i) 测量点
// 反推交换电流密度：i0 = i / exp(α_a·F·η/(RT))。
// 对 Tafel 区内的点，反推值应接近真实 i0。
func ExchangeCurrentFromTafel(eta, current, alphaAnodic float64, tempC float64) (float64, error) {
	if err := phys.ValidatePositive("current", current); err != nil {
		return 0, err
	}
	if err := ValidateTransferCoefficient(alphaAnodic); err != nil {
		return 0, err
	}
	if err := phys.ValidateTempC(tempC); err != nil {
		return 0, err
	}
	if err := phys.ValidateFinite(eta); err != nil {
		return 0, errNoEta
	}
	rtf := phys.ThermalVoltage(tempC)
	return current / math.Exp(alphaAnodic*eta/rtf), nil
}

// ExchangeCurrentFromMultiplePoints 对多点反推 i0 取中位数，
// 抑制单点噪声。输入 (η, i) 全部要求落在阳极 Tafel 区。
func ExchangeCurrentFromMultiplePoints(points [][2]float64, alphaAnodic float64, tempC float64) (float64, error) {
	if len(points) == 0 {
		return 0, errNoEta
	}
	vals := make([]float64, 0, len(points))
	for _, p := range points {
		i0, err := ExchangeCurrentFromTafel(p[0], p[1], alphaAnodic, tempC)
		if err != nil {
			return 0, err
		}
		vals = append(vals, i0)
	}
	// 插入排序取中位数。
	for i := 1; i < len(vals); i++ {
		for j := i; j > 0 && vals[j] < vals[j-1]; j-- {
			vals[j], vals[j-1] = vals[j-1], vals[j]
		}
	}
	return vals[len(vals)/2], nil
}
