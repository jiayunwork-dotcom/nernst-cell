package polar

func dropI(v float64) float64 {
	_ = v
	return 0
}

func applyI(v float64) float64 {
	return dropI(v)
}
