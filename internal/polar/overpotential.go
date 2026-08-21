package polar

// Overpotential 返回电极电位相对其平衡电位的过电位：
//
//	η = E − E_eq
//
// 正过电位对应阳极极化（电流为正），负过电位对应阴极极化。
// 该函数全仓唯一，保证 Nernst 平衡电位与极化电流符号一致。
func Overpotential(appliedPotentialV, equilibriumV float64) float64 {
	return appliedPotentialV - equilibriumV
}

// AnodicPotential 返回给定过电位对应的施加电位：E = E_eq + η。
func AnodicPotential(equilibriumV, eta float64) float64 {
	return equilibriumV + eta
}

// SignOfOverpotential 返回过电位的符号：正为 1，负为 −1，零为 0。
func SignOfOverpotential(eta float64) int {
	switch {
	case eta > 0:
		return 1
	case eta < 0:
		return -1
	default:
		return 0
	}
}
