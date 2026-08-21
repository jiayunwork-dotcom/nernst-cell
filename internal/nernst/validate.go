package nernst

import (
	"errors"

	"nernst-cell/internal/phys"
)

// Validate 校验 Nernst 输入：
// n 必须为正整数；T（摄氏）必须为正；活度必须为正且有限；
// 标准电位必须有限。任一不满足即返回 error。
func Validate(in Input) error {
	if in.Electrons <= 0 {
		return errors.New("electrons must be a positive integer")
	}
	if err := phys.ValidateTempC(in.TemperatureC); err != nil {
		return err
	}
	if err := phys.ValidatePositive("ox activity", in.OxActivity); err != nil {
		return err
	}
	if err := phys.ValidatePositive("red activity", in.RedActivity); err != nil {
		return err
	}
	if err := phys.ValidateFinite(in.StandardPotentialV); err != nil {
		return errors.New("standard potential must be finite")
	}
	return nil
}
