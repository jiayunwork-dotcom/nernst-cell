package polar

import (
	"errors"
	"math"

	"nernst-cell/internal/phys"
)

// errOverpotentialNotFinite 在过电位非有限时返回。
var errOverpotentialNotFinite = errors.New("overpotential must be finite")

// errNoEta 是过电位参数缺失/非法时的通用错误。
var errNoEta = errors.New("overpotential must be finite")

// ValidateKinetics 校验极化动力学公共输入：交换电流密度必须为
// 正且有限，传递系数必须在 (0, n) 内，温度必须合法。
func ValidateKinetics(i0, alphaAnodic float64, electrons int, tempC float64) error {
	if err := phys.ValidatePositive("exchange current density", i0); err != nil {
		return err
	}
	if err := ValidateTransferCoefficient(alphaAnodic); err != nil {
		return err
	}
	if electrons <= 0 {
		return errors.New("electrons must be a positive integer")
	}
	if alphaAnodic >= float64(electrons) {
		return errors.New("alpha_anodic must be less than electrons")
	}
	return phys.ValidateTempC(tempC)
}

// ValidateTransferCoefficient 校验单个传递系数：必须在 (0, n) 内，
// 但此处只做与 n 无关的正有限检查；n 相关的上限由调用方完成。
func ValidateTransferCoefficient(alpha float64) error {
	if math.IsNaN(alpha) || math.IsInf(alpha, 0) {
		return errors.New("transfer coefficient must be finite")
	}
	if alpha <= 0 {
		return errors.New("transfer coefficient must be positive")
	}
	return nil
}

// ValidateIVInput 校验极化曲线输入：动力学参数、过电位网格
// 与 Nernst 输入全部合法才放行。
func ValidateIVInput(in IVInput) error {
	if err := ValidateKinetics(in.ExchangeCurrentDensity, in.AlphaAnodic, in.Nernst.Electrons, in.Nernst.TemperatureC); err != nil {
		return err
	}
	if err := phys.ValidateFinite(in.EtaMinV); err != nil {
		return errors.New("eta_min must be finite")
	}
	if err := phys.ValidateFinite(in.EtaMaxV); err != nil {
		return errors.New("eta_max must be finite")
	}
	if in.EtaMinV >= in.EtaMaxV {
		return errors.New("eta_min must be less than eta_max")
	}
	if in.EtaPoints < 2 {
		return errors.New("eta_points must be at least 2")
	}
	if in.EtaPoints > 500 {
		return errors.New("eta_points must be at most 500")
	}
	return nil
}
