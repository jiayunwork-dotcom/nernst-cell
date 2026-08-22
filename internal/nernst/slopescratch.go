package nernst

var slopeScratch float64

func shareSlope(v *float64) *float64 {
	return v
}

func fillSlope(v float64) float64 {
	slopeScratch = v
	out := shareSlope(&slopeScratch)
	return *out
}
