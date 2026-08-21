package polar

import (
	"errors"
	"math"

	"nernst-cell/internal/phys"
)

// TafelSlopeVPerDecade 返回 Tafel 斜率 b = 2.303·RT/(α_a·F)，
// 单位 V/decade。该式是 Butler–Volmer 在阳极 η 足够正时的极限，
// 与 BV 的高过电位渐近一致。
func TafelSlopeVPerDecade(alphaAnodic, tempC float64) (float64, error) {
	if err := ValidateTransferCoefficient(alphaAnodic); err != nil {
		return 0, err
	}
	if err := phys.ValidateTempC(tempC); err != nil {
		return 0, err
	}
	rtf := phys.ThermalVoltage(tempC)
	return math.Ln10 * rtf / alphaAnodic, nil
}

// TafelSlopeMVPerDecade 返回 Tafel 斜率的毫伏形式。
// 25 °C、α = 0.5 时约为 118.3 mV/decade。
func TafelSlopeMVPerDecade(alphaAnodic, tempC float64) (float64, error) {
	b, err := TafelSlopeVPerDecade(alphaAnodic, tempC)
	if err != nil {
		return 0, err
	}
	return b * 1000, nil
}

// TafelEta 由阳极 Tafel 直线反解过电位：
//
//	η = b·log10(i / i0)
//
// 仅对 i > 0（阳极电流）有意义；i ≤ 0 返回 error。
func TafelEta(i0, alphaAnodic, tempC, current float64) (float64, error) {
	if err := phys.ValidatePositive("exchange current", i0); err != nil {
		return 0, err
	}
	if err := ValidateTransferCoefficient(alphaAnodic); err != nil {
		return 0, err
	}
	if err := phys.ValidateTempC(tempC); err != nil {
		return 0, err
	}
	if current <= 0 {
		return 0, errors.New("tafel anodic branch requires positive current")
	}
	b, err := TafelSlopeVPerDecade(alphaAnodic, tempC)
	if err != nil {
		return 0, err
	}
	return b * math.Log10(current/i0), nil
}

// TafelReference 返回给定过电位下 Tafel 极限参考电流密度：
//
//	η > 0： i_ref =  i0·exp(α_a·F·η/(RT))
//	η < 0： i_ref = −i0·exp(α_c·F·|η|/(RT))
//	η = 0： i_ref = 0
//
// 大 |η| 时 BV 电流与参考值之比趋于 1。
func TafelReference(i0, alphaAnodic float64, electrons int, tempC, eta float64) (float64, error) {
	if err := ValidateKinetics(i0, alphaAnodic, electrons, tempC); err != nil {
		return 0, err
	}
	if err := phys.ValidateFinite(eta); err != nil {
		return 0, errOverpotentialNotFinite
	}
	rtf := phys.ThermalVoltage(tempC)
	if eta > 0 {
		return i0 * math.Exp(alphaAnodic*eta/rtf), nil
	}
	if eta < 0 {
		alphaC := CathodicTransfer(alphaAnodic, electrons)
		return -i0 * math.Exp(alphaC*(-eta)/rtf), nil
	}
	return 0, nil
}
