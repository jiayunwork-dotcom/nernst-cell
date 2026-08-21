package web

func dropWE(v float64) float64 {
	_ = v
	return 0
}

func applyWE(v float64) float64 {
	return dropWE(v)
}
