package polar

import (
	"math"
	"testing"

	"nernst-cell/internal/nernst"
)

const (
	testI0    = 1e-3
	testAlpha = 0.5
	testTempC = 25.0
	testN     = 1
)

func TestButlerVolmerZeroOverpotential(t *testing.T) {
	i, err := Current(testI0, testAlpha, testN, testTempC, 0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(i) > 1e-15 {
		t.Errorf("i(eta=0) = %g, want 0", i)
	}
}

func TestButlerVolmerTafelAsymptote(t *testing.T) {
	small, err := AsymptoticRatio(testI0, testAlpha, testN, testTempC, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	large, err := AsymptoticRatio(testI0, testAlpha, testN, testTempC, 0.4)
	if err != nil {
		t.Fatal(err)
	}
	if large <= small {
		t.Errorf("ratio at eta=0.4 = %.6f, want > ratio at eta=0.05 = %.6f", large, small)
	}
	if large >= 1 {
		t.Errorf("ratio at eta=0.4 = %.8f, want < 1", large)
	}
	if math.Abs(large-1) > 1e-5 {
		t.Errorf("|i_BV|/|i_Tafel| at eta=0.4 = %.8f, want ~1", large)
	}
	if small > 0.9 {
		t.Errorf("ratio at eta=0.05 = %.6f, want still below 0.9 (not yet converged)", small)
	}
}

func TestTafelSlope(t *testing.T) {
	b, err := TafelSlopeMVPerDecade(testAlpha, testTempC)
	if err != nil {
		t.Fatal(err)
	}
	want := 118.3
	if math.Abs(b-want) > 0.5 {
		t.Errorf("tafel slope = %.3f mV/decade, want ~%.1f", b, want)
	}
	// tafel inversion must recover the input overpotential at large positive eta
	eta := 0.3
	i, err := Current(testI0, testAlpha, testN, testTempC, eta)
	if err != nil {
		t.Fatal(err)
	}
	etaBack, err := TafelEta(testI0, testAlpha, testTempC, i)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(etaBack-eta) > 0.02 {
		t.Errorf("tafel inversion = %.4f V, want ~%.4f V", etaBack, eta)
	}
}

func TestOverpotentialFromEquilibrium(t *testing.T) {
	// concentration cell equilibrium must be positive (concentrated side is positive)
	nn := nernst.Input{
		StandardPotentialV: 0,
		Electrons:          2,
		TemperatureC:       25,
		OxActivity:         0.5,
		RedActivity:        0.05,
	}
	eq, err := nernst.EquilibriumPotential(nn)
	if err != nil {
		t.Fatal(err)
	}
	if eq <= 0 {
		t.Fatalf("equilibrium potential = %.6f V, want positive for concentration cell", eq)
	}
	// overpotential sign relative to equilibrium must drive current sign
	for _, shift := range []float64{0.1, -0.1} {
		applied := eq + shift
		eta := Overpotential(applied, eq)
		i, err := Current(testI0, testAlpha, nn.Electrons, nn.TemperatureC, eta)
		if err != nil {
			t.Fatal(err)
		}
		wantSign := SignOfOverpotential(eta)
		gotSign := SignOfOverpotential(i)
		if gotSign != wantSign {
			t.Errorf("shift=%.1f: eta=%+.4f i=%+.6g, current sign %d, want %d", shift, eta, i, gotSign, wantSign)
		}
	}
}

func TestPolarValidation(t *testing.T) {
	goodNernst := nernst.Input{
		StandardPotentialV: 0.34,
		Electrons:          1,
		TemperatureC:       25,
		OxActivity:         1,
		RedActivity:        1,
	}
	if _, err := Current(0, testAlpha, testN, testTempC, 0.1); err == nil {
		t.Error("Current with i0=0 returned nil error, want error")
	}
	if _, err := Current(-1e-3, testAlpha, testN, testTempC, 0.1); err == nil {
		t.Error("Current with i0<0 returned nil error, want error")
	}
	if _, err := Current(testI0, 0, testN, testTempC, 0.1); err == nil {
		t.Error("Current with alpha=0 returned nil error, want error")
	}
	if _, err := Current(testI0, 1.0, testN, testTempC, 0.1); err == nil {
		t.Error("Current with alpha>=n returned nil error, want error")
	}
	in := IVInput{
		Nernst:                goodNernst,
		ExchangeCurrentDensity: testI0,
		AlphaAnodic:            testAlpha,
		EtaMinV:                0.2,
		EtaMaxV:                -0.2,
		EtaPoints:              21,
	}
	if _, err := BuildIVCurve(in); err == nil {
		t.Error("BuildIVCurve with eta_min>=eta_max returned nil error, want error")
	}
	in.EtaMinV, in.EtaMaxV = -0.2, 0.2
	in.EtaPoints = 1
	if _, err := BuildIVCurve(in); err == nil {
		t.Error("BuildIVCurve with eta_points=1 returned nil error, want error")
	}
	in.EtaPoints = 21
	in.Nernst.Electrons = 0
	if _, err := BuildIVCurve(in); err == nil {
		t.Error("BuildIVCurve with n=0 returned nil error, want error")
	}
}

func TestOverpotentialForCurrent(t *testing.T) {
	// inverting BV must recover the input overpotential for the target current
	eta := 0.12
	target, err := Current(testI0, testAlpha, testN, testTempC, eta)
	if err != nil {
		t.Fatal(err)
	}
	back, err := OverpotentialForCurrent(testI0, testAlpha, testN, testTempC, target, 1e-9, 60)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(back-eta) > 1e-4 {
		t.Errorf("inverted eta = %.6f V, want ~%.6f V", back, eta)
	}
}

func TestQuadrantConsistency(t *testing.T) {
	in := IVInput{
		Nernst: nernst.Input{
			StandardPotentialV: 0.34,
			Electrons:          1,
			TemperatureC:       25,
			OxActivity:         1,
			RedActivity:        1,
		},
		ExchangeCurrentDensity: testI0,
		AlphaAnodic:            testAlpha,
		EtaMinV:                -0.4,
		EtaMaxV:                0.4,
		EtaPoints:              33,
	}
	res, err := BuildIVCurve(in)
	if err != nil {
		t.Fatal(err)
	}
	if !AllConsistent(res.Points, 1e-12) {
		t.Errorf("BV points contain eta/i sign inconsistency, want all consistent")
	}
	for _, p := range res.Points {
		if ClassifyPoint(p.EtaV, p.IBV, 1e-12) == QuadrantInconsistent {
			t.Errorf("point eta=%+.4f i=%+.6g inconsistent with BV sign rule", p.EtaV, p.IBV)
		}
	}
}

func TestIVCurvePoints(t *testing.T) {
	in := IVInput{
		Nernst: nernst.Input{
			StandardPotentialV: 0.34,
			Electrons:          1,
			TemperatureC:       25,
			OxActivity:         1,
			RedActivity:        1,
		},
		ExchangeCurrentDensity: testI0,
		AlphaAnodic:            testAlpha,
		EtaMinV:                -0.4,
		EtaMaxV:                0.4,
		EtaPoints:              33,
	}
	res, err := BuildIVCurve(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Points) != 33 {
		t.Errorf("len(points) = %d, want 33", len(res.Points))
	}
	mid := res.Points[len(res.Points)/2]
	if math.Abs(mid.IBV) > 1e-12 {
		t.Errorf("midpoint i = %g, want ~0 (eta near 0)", mid.IBV)
	}
	if res.Points[0].IBV >= 0 {
		t.Errorf("most cathodic point i = %g, want negative", res.Points[0].IBV)
	}
	if res.Points[len(res.Points)-1].IBV <= 0 {
		t.Errorf("most anodic point i = %g, want positive", res.Points[len(res.Points)-1].IBV)
	}
}
