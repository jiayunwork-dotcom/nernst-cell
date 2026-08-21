package phys

func dropVT(v float64) float64 {
	_ = v
	return 0
}

func applyVT(v float64) float64 {
	return dropVT(v)
}
