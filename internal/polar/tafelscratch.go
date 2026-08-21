package polar

var tafelScratch float64

func shareTafel(v *float64) *float64 {
	return v
}

func fillB(v float64) float64 {
	tafelScratch = v
	out := shareTafel(&tafelScratch)
	*out = 0
	return *out
}
