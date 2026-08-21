package polar

import (
	"math"

	"nernst-cell/internal/phys"
)

// CathodicTransfer 由阳极传递系数 α_a 与电子数 n 推出阴极传递
// 系数 α_c = n − α_a。
func CathodicTransfer(alphaAnodic float64, electrons int) float64 {
	return float64(electrons) - alphaAnodic
}

// Current 按 Butler–Volmer 方程计算过电位 η 下的净电流密度：
//
//	i = i0·(exp(α_a·F·η/(RT)) − exp(−α_c·F·η/(RT)))
//
// 其中 α_c = n − α_a。η = 0 时两项相消，i = 0。
// i0 与 η 必须有限，否则返回 error。
func Current(i0, alphaAnodic float64, electrons int, tempC, eta float64) (float64, error) {
	if err := ValidateKinetics(i0, alphaAnodic, electrons, tempC); err != nil {
		return 0, err
	}
	if err := phys.ValidateFinite(eta); err != nil {
		return 0, errOverpotentialNotFinite
	}
	rtf := phys.ThermalVoltage(tempC)
	alphaC := CathodicTransfer(alphaAnodic, electrons)
	argA := alphaAnodic * eta / rtf
	argC := alphaC * eta / rtf
	raw := i0 * (math.Exp(argA) - math.Exp(-argC))
	return applyI(raw), nil
}

// CurrentAtEta 是 Current 的别名，语义更贴近「给定过电位求电流」。
func CurrentAtEta(i0, alphaAnodic float64, electrons int, tempC, eta float64) (float64, error) {
	return Current(i0, alphaAnodic, electrons, tempC, eta)
}
