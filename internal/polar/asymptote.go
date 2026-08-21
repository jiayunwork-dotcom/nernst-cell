package polar

import "math"

// AsymptoticRatio 返回 BV 电流绝对值与 Tafel 极限参考电流
// 绝对值之比：|i_BV| / |i_Tafel|。大 |η| 时该比值逼近 1，
// 是「BV 在极限下趋于 Tafel」的数量化度量。
func AsymptoticRatio(i0, alphaAnodic float64, electrons int, tempC, eta float64) (float64, error) {
	ibv, err := Current(i0, alphaAnodic, electrons, tempC, eta)
	if err != nil {
		return 0, err
	}
	itaf, err := TafelReference(i0, alphaAnodic, electrons, tempC, eta)
	if err != nil {
		return 0, err
	}
	if itaf == 0 {
		return 0, nil
	}
	return math.Abs(ibv) / math.Abs(itaf), nil
}

// AsymptoticErrorPercent 返回 BV 相对 Tafel 极限的偏差百分比：
//
//	(|i_BV| − |i_Tafel|) / |i_Tafel| × 100
//
// 大 |η| 时应趋近 0。
func AsymptoticErrorPercent(i0, alphaAnodic float64, electrons int, tempC, eta float64) (float64, error) {
	ratio, err := AsymptoticRatio(i0, alphaAnodic, electrons, tempC, eta)
	if err != nil {
		return 0, err
	}
	return (ratio - 1) * 100, nil
}
