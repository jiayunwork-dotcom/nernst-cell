package web

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const cuConcExample = `{
  "standard_potential_v": 0.0,
  "electrons": 2,
  "temperature_c": 25,
  "ox_activity": 0.5,
  "red_activity": 0.05
}`

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	return NewServer(Assets{
		Examples: map[string][]byte{
			"cu-conc": []byte(cuConcExample),
		},
	})
}

func post(t *testing.T, handler http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestWebNernstEndpoint(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/nernst", cuConcExample)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		E                 float64 `json:"equilibrium_potential_v"`
		SlopeMVPerDecade  float64 `json:"slope_mv_per_decade"`
		DecadeShiftV      float64 `json:"decade_shift_v"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	// cu concentration cell: n=2, ratio 10, E = 59.16/2 ≈ 29.58 mV, positive
	want := 0.02958
	if math.Abs(resp.E-want) > 1e-4 {
		t.Errorf("E = %.6f V, want ~%.5f V", resp.E, want)
	}
	if resp.E <= 0 {
		t.Errorf("E = %.6f V, want positive (concentrated side is positive)", resp.E)
	}
	if math.Abs(resp.SlopeMVPerDecade-29.58) > 0.01 {
		t.Errorf("slope = %.3f mV/decade, want ~29.58", resp.SlopeMVPerDecade)
	}
}

func TestWebNernstInvalidElectrons(t *testing.T) {
	handler := testHandler(t)
	body := `{"standard_potential_v":0,"electrons":0,"temperature_c":25,"ox_activity":1,"red_activity":1}`
	rec := post(t, handler, "/api/nernst", body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "electrons") {
		t.Errorf("error = %q, want it to mention electrons", resp.Error)
	}
}

func TestWebNernstInvalidActivity(t *testing.T) {
	handler := testHandler(t)
	body := `{"standard_potential_v":0,"electrons":1,"temperature_c":25,"ox_activity":-1,"red_activity":1}`
	rec := post(t, handler, "/api/nernst", body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "positive") {
		t.Errorf("error = %q, want it to mention positive", resp.Error)
	}
}

func TestWebIVEndpoint(t *testing.T) {
	handler := testHandler(t)
	body := `{
		"standard_potential_v":0.34,"electrons":1,"temperature_c":25,
		"ox_activity":1,"red_activity":1,
		"exchange_current_density":1e-3,"alpha":0.5,
		"eta_min_v":-0.4,"eta_max_v":0.4,"eta_points":33
	}`
	rec := post(t, handler, "/api/iv", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		SlopeMVPerDecade float64 `json:"slope_mv_per_decade"`
		TafelSlope       float64 `json:"tafel_slope_mv_per_decade"`
		Points           []struct {
			EtaV float64 `json:"eta_v"`
			IBV  float64 `json:"i_bv"`
		} `json:"points"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Points) != 33 {
		t.Errorf("len(points) = %d, want 33", len(resp.Points))
	}
	if math.Abs(resp.SlopeMVPerDecade-59.16) > 0.01 {
		t.Errorf("nernst slope = %.3f mV/decade, want ~59.16", resp.SlopeMVPerDecade)
	}
	if math.Abs(resp.TafelSlope-118.3) > 0.5 {
		t.Errorf("tafel slope = %.3f mV/decade, want ~118.3", resp.TafelSlope)
	}
	mid := resp.Points[len(resp.Points)/2]
	if math.Abs(mid.IBV) > 1e-12 {
		t.Errorf("midpoint i = %g, want ~0 (eta near 0)", mid.IBV)
	}
	if resp.Points[0].IBV >= 0 {
		t.Errorf("most cathodic point i = %g, want negative", resp.Points[0].IBV)
	}
}

func TestWebIVInvalidExchangeCurrent(t *testing.T) {
	handler := testHandler(t)
	body := `{
		"standard_potential_v":0.34,"electrons":1,"temperature_c":25,
		"ox_activity":1,"red_activity":1,
		"exchange_current_density":0,"alpha":0.5,
		"eta_min_v":-0.2,"eta_max_v":0.2,"eta_points":21
	}`
	rec := post(t, handler, "/api/iv", body)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "exchange current") {
		t.Errorf("error = %q, want it to mention exchange current", resp.Error)
	}
}

func TestWebExamplesEndpoint(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/examples", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp struct {
		CuConc json.RawMessage `json:"cu-conc"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.CuConc) == 0 {
		t.Errorf("examples payload missing cu-conc")
	}
	var ex nernstRequest
	if err := json.Unmarshal(resp.CuConc, &ex); err != nil {
		t.Fatal(err)
	}
	if ex.Electrons != 2 || ex.OxActivity != 0.5 {
		t.Errorf("cu-conc example parsed as electrons=%v ox=%v, want 2 and 0.5", ex.Electrons, ex.OxActivity)
	}
}
