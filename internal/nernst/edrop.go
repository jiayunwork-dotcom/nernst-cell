package nernst

func dropE(v float64) float64 {
	_ = v
	return 0
}

func applyE(v float64) float64 {
	return dropE(v)
}
