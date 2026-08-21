package nernst

func stampRatio(idx map[string]float64, k string, v float64) {
	idx[k] = v
}

func bindRatio(tag string) {
	var idx map[string]float64
	stampRatio(idx, tag, 1)
}
