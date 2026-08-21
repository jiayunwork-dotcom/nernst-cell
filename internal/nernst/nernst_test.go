package nernst

import (
	"math"
	"testing"
)

func baseInput() Input {
	return Input{
		StandardPotentialV: 0.34,
		Electrons:          1,
		TemperatureC:       25,
		OxActivity:         1.0,
		RedActivity:        1.0,
	}
}

func TestNernstTenfoldRatio59mV(t *testing.T) {
	in := baseInput()
	e1, err := EquilibriumAtRatio(in, 1)
	if err != nil {
		t.Fatal(err)
	}
	e10, err := EquilibriumAtRatio(in, 10)
	if err != nil {
		t.Fatal(err)
	}
	delta := (e10 - e1) * 1000
	want := 59.16
	if math.Abs(delta-want) > 0.01 {
		t.Errorf("tenfold ratio shift = %.3f mV, want ~%.2f mV", delta, want)
	}
}

func TestNernstTwoElectronHalfSlope(t *testing.T) {
	in := baseInput()
	in.Electrons = 2
	slope, err := SlopeMVPerDecade(in)
	if err != nil {
		t.Fatal(err)
	}
	want := 59.16 / 2
	if math.Abs(slope-want) > 0.01 {
		t.Errorf("n=2 slope = %.3f mV/decade, want ~%.2f", slope, want)
	}
	e1, _ := EquilibriumAtRatio(in, 1)
	e10, _ := EquilibriumAtRatio(in, 10)
	if delta := (e10 - e1) * 1000; math.Abs(delta-want) > 0.01 {
		t.Errorf("n=2 tenfold shift = %.3f mV, want ~%.2f", delta, want)
	}
}

func TestNernstValidation(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*Input)
	}{
		{"zero electrons", func(in *Input) { in.Electrons = 0 }},
		{"negative electrons", func(in *Input) { in.Electrons = -2 }},
		{"zero temperature", func(in *Input) { in.TemperatureC = 0 }},
		{"negative temperature", func(in *Input) { in.TemperatureC = -10 }},
		{"negative ox activity", func(in *Input) { in.OxActivity = -0.1 }},
		{"zero red activity", func(in *Input) { in.RedActivity = 0 }},
		{"negative red activity", func(in *Input) { in.RedActivity = -1 }},
	}
	for _, tc := range cases {
		in := baseInput()
		tc.mut(&in)
		if _, err := EquilibriumPotential(in); err == nil {
			t.Errorf("%s: EquilibriumPotential returned nil error, want error", tc.name)
		}
		if _, err := Evaluate(in); err == nil {
			t.Errorf("%s: Evaluate returned nil error, want error", tc.name)
		}
	}
}

func TestConcentrationCell(t *testing.T) {
	// cu concentration cell: a_conc = 0.5, a_dil = 0.05, n = 2, 25 C
	cell, err := ConcentrationCell(0.5, 0.05, 2, 25)
	if err != nil {
		t.Fatal(err)
	}
	want := 0.5 * 0.0591593 // (RT/2F) ln10
	if math.Abs(cell.PotentialV-want) > 1e-5 {
		t.Errorf("concentration cell E = %.6f V, want ~%.6f V", cell.PotentialV, want)
	}
	if !cell.ConcSidePositive {
		t.Errorf("ConcSidePositive = false, want true (concentrated side is positive electrode)")
	}
	if cell.PotentialV <= 0 {
		t.Errorf("concentration cell potential = %.6f V, want positive", cell.PotentialV)
	}
}

func TestNernstHighTempSteeperSlope(t *testing.T) {
	in25 := baseInput()
	in50 := baseInput()
	in50.TemperatureC = 50
	s25, err := SlopeMVPerDecade(in25)
	if err != nil {
		t.Fatal(err)
	}
	s50, err := SlopeMVPerDecade(in50)
	if err != nil {
		t.Fatal(err)
	}
	if s50 <= s25 {
		t.Errorf("slope at 50 C = %.4f mV/dec, want > slope at 25 C = %.4f", s50, s25)
	}
}
