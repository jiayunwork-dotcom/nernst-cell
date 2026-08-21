package phys

import (
	"math"
	"testing"
)

func TestThermalVoltageStandard(t *testing.T) {
	got := ThermalVoltageStandard()
	want := 0.0256926
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("RT/F at 25 C = %.7f, want ~%.7f", got, want)
	}
	decade := DecadeVoltageStandard()
	wantMV := 59.16e-3
	if math.Abs(decade-wantMV) > 1e-4 {
		t.Errorf("RT/F*ln10 at 25 C = %.6f V, want ~%.6f V", decade, wantMV)
	}
}

func TestTempCToKAndValidation(t *testing.T) {
	if got := TempCToK(25); math.Abs(got-298.15) > 1e-12 {
		t.Errorf("TempCToK(25) = %.6f, want 298.15", got)
	}
	if got := KToTempC(TempCToK(80)); math.Abs(got-80) > 1e-9 {
		t.Errorf("round trip KToTempC(TempCToK(80)) = %.6f, want 80", got)
	}
	for _, bad := range []float64{0, -25, math.NaN(), math.Inf(1)} {
		if err := ValidateTempC(bad); err == nil {
			t.Errorf("ValidateTempC(%v) = nil, want error", bad)
		}
	}
	if err := ValidateTempC(25); err != nil {
		t.Errorf("ValidateTempC(25) = %v, want nil", err)
	}
}

func TestElectronCountValidation(t *testing.T) {
	for _, bad := range []float64{0, -1, 2.5, math.NaN(), math.Inf(1)} {
		if err := ValidateElectronCount(bad); err == nil {
			t.Errorf("ValidateElectronCount(%v) = nil, want error", bad)
		}
	}
	for _, good := range []float64{1, 2, 3} {
		if err := ValidateElectronCount(good); err != nil {
			t.Errorf("ValidateElectronCount(%v) = %v, want nil", good, err)
		}
	}
}
