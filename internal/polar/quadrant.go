package polar

// Quadrant 描述 (η, i) 点所在的电流–过电位象限。
type Quadrant int

const (
	// QuadrantAnodic 是阳极区：η > 0 且 i > 0。
	QuadrantAnodic Quadrant = iota
	// QuadrantCathodic 是阴极区：η < 0 且 i < 0。
	QuadrantCathodic
	// QuadrantZero 是平衡点：η ≈ 0 且 i ≈ 0。
	QuadrantZero
	// QuadrantInconsistent 是符号矛盾点：η 与 i 异号。
	QuadrantInconsistent
)

// ClassifyPoint 按符号把 (η, i) 归入象限。
// 交叉规则要求 BV 电流与过电位同号：η > 0 必阳极电流为正，
// η < 0 必阴极电流为负；违背即 Inconsistent。
func ClassifyPoint(eta, i float64, tol float64) Quadrant {
	both := func(a, b float64) bool {
		return a > tol && b > tol || a < -tol && b < -tol
	}
	if both(eta, i) {
		if eta > 0 {
			return QuadrantAnodic
		}
		return QuadrantCathodic
	}
	if absEta(eta) <= tol && absEta(i) <= tol {
		return QuadrantZero
	}
	return QuadrantInconsistent
}

// AllConsistent 检查点列是否全部满足「η 与 i 同号」的 BV 契约。
// 除平衡点外，任何异号点都判定为不一致。
func AllConsistent(points []IVPoint, tol float64) bool {
	for _, p := range points {
		if ClassifyPoint(p.EtaV, p.IBV, tol) == QuadrantInconsistent {
			return false
		}
	}
	return true
}

func absEta(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
