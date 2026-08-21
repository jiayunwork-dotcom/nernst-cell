package polar

import (
	"errors"
	"math"

	"nernst-cell/internal/phys"
)

// OverpotentialForCurrent 对 Butler–Volmer 求逆：给定目标电流
// （正为阳极），二分搜索满足 Current(η) = targetI 的过电位 η。
// 迭代最多 maxIter 次，容差 tol（A/cm²）。BV 对 η 单调，解唯一。
func OverpotentialForCurrent(i0, alphaAnodic float64, electrons int, tempC, targetI float64, tol float64, maxIter int) (float64, error) {
	if err := phys.ValidatePositive("target current", targetI); err != nil {
		return 0, err
	}
	if tol <= 0 {
		return 0, errors.New("tolerance must be positive")
	}
	if maxIter <= 0 {
		return 0, errors.New("max iterations must be positive")
	}
	// 扩展上界直至 Current 超过目标。
	lo := 0.0
	hi := 0.05
	for i := 0; i < 40; i++ {
		v, err := Current(i0, alphaAnodic, electrons, tempC, hi)
		if err != nil {
			return 0, err
		}
		if v >= targetI {
			break
		}
		hi *= 2
	}
	for it := 0; it < maxIter; it++ {
		mid := (lo + hi) / 2
		v, err := Current(i0, alphaAnodic, electrons, tempC, mid)
		if err != nil {
			return 0, err
		}
		if math.Abs(v-targetI) <= tol {
			return mid, nil
		}
		if v < targetI {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2, nil
}

// HalfWaveOverpotential 返回对称系数下净电流等于 i0 时的过电位
// η₁/₂，即 BV 中 exp(α_a Fη/RT) − exp(−α_c Fη/RT) = 1 的解。
// 对 α = 0.5、n = 1 有解析式 η = (RT/F)·asinh(1)/2 ≈ 0.0113 V。
func HalfWaveOverpotential(alphaAnodic float64, electrons int, tempC float64) (float64, error) {
	if err := phys.ValidateTempC(tempC); err != nil {
		return 0, err
	}
	if err := ValidateTransferCoefficient(alphaAnodic); err != nil {
		return 0, err
	}
	return OverpotentialForCurrent(1, alphaAnodic, electrons, tempC, 1.0, 1e-12, 60)
}
