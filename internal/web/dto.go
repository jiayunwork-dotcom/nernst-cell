package web

import (
	"nernst-cell/internal/nernst"
	"nernst-cell/internal/phys"
	"nernst-cell/internal/polar"
)

// nernstRequest 是 /api/nernst 的请求体。
type nernstRequest struct {
	StandardPotentialV float64 `json:"standard_potential_v"`
	Electrons          float64 `json:"electrons"`
	TemperatureC       float64 `json:"temperature_c"`
	OxActivity         float64 `json:"ox_activity"`
	RedActivity        float64 `json:"red_activity"`
}

// nernstResponse 是 /api/nernst 的响应体。
type nernstResponse struct {
	EquilibriumPotentialV float64 `json:"equilibrium_potential_v"`
	SlopeMVPerDecade      float64 `json:"slope_mv_per_decade"`
	ThermalVoltageV       float64 `json:"thermal_voltage_v"`
	DecadeShiftV          float64 `json:"decade_shift_v"`
	ActivityRatio         float64 `json:"activity_ratio"`
	TemperatureC          float64 `json:"temperature_c"`
	Electrons             int     `json:"electrons"`
}

// toInput 把请求体转成 nernst.Input，先做字段级校验。
func (r nernstRequest) toInput() (nernst.Input, error) {
	if err := phys.ValidateElectronCount(r.Electrons); err != nil {
		return nernst.Input{}, err
	}
	return nernst.Input{
		StandardPotentialV: r.StandardPotentialV,
		Electrons:          int(r.Electrons),
		TemperatureC:       r.TemperatureC,
		OxActivity:         r.OxActivity,
		RedActivity:        r.RedActivity,
	}, nil
}

// toIVInput 把请求体转成 polar.IVInput，校验电子数。
// alpha 缺省（0）时按 α = 0.5 处理。
func (r ivRequest) toIVInput() (polar.IVInput, error) {
	if err := phys.ValidateElectronCount(r.Electrons); err != nil {
		return polar.IVInput{}, err
	}
	alpha := r.Alpha
	if alpha == 0 {
		alpha = 0.5
	}
	return polar.IVInput{
		Nernst: nernst.Input{
			StandardPotentialV: r.StandardPotentialV,
			Electrons:          int(r.Electrons),
			TemperatureC:       r.TemperatureC,
			OxActivity:         r.OxActivity,
			RedActivity:        r.RedActivity,
		},
		ExchangeCurrentDensity: r.ExchangeCurrentDensity,
		AlphaAnodic:            alpha,
		EtaMinV:                r.EtaMinV,
		EtaMaxV:                r.EtaMaxV,
		EtaPoints:              r.EtaPoints,
	}, nil
}

// ivRequest 是 /api/iv 的请求体。
type ivRequest struct {
	StandardPotentialV    float64 `json:"standard_potential_v"`
	Electrons             float64 `json:"electrons"`
	TemperatureC          float64 `json:"temperature_c"`
	OxActivity            float64 `json:"ox_activity"`
	RedActivity           float64 `json:"red_activity"`
	ExchangeCurrentDensity float64 `json:"exchange_current_density"`
	Alpha                 float64 `json:"alpha"`
	EtaMinV               float64 `json:"eta_min_v"`
	EtaMaxV               float64 `json:"eta_max_v"`
	EtaPoints             int     `json:"eta_points"`
}

// ivPointDTO 是响应中单点的 JSON 表示。
type ivPointDTO struct {
	EtaV   float64 `json:"eta_v"`
	IBV    float64 `json:"i_bv"`
	ITafel float64 `json:"i_tafel"`
}

// ivResponse 是 /api/iv 的响应体。
type ivResponse struct {
	EquilibriumPotentialV   float64      `json:"equilibrium_potential_v"`
	SlopeMVPerDecade        float64      `json:"slope_mv_per_decade"`
	TafelSlopeMVPerDecade   float64      `json:"tafel_slope_mv_per_decade"`
	ExchangeCurrentDensity  float64      `json:"exchange_current_density"`
	AlphaAnodic             float64      `json:"alpha_anodic"`
	AlphaCathodic           float64      `json:"alpha_cathodic"`
	Points                  []ivPointDTO `json:"points"`
}

// fromResult 把 nernst.Result 转成响应体，数值统一舍入。
func fromResult(res nernst.Result, electrons int) nernstResponse {
	return nernstResponse{
		EquilibriumPotentialV: roundE(applyWE(res.EquilibriumPotentialV)),
		SlopeMVPerDecade:      roundSlope(res.SlopeMVPerDecade),
		ThermalVoltageV:       roundE(res.ThermalVoltageV),
		DecadeShiftV:          roundE(res.DecadeShiftV),
		ActivityRatio:         roundSlope(res.ActivityRatio),
		TemperatureC:          roundSlope(res.TemperatureC),
		Electrons:             electrons,
	}
}

// fromIVResult 把 polar.IVResult 转成响应体，数值统一舍入。
func fromIVResult(res polar.IVResult) ivResponse {
	points := make([]ivPointDTO, 0, len(res.Points))
	for _, p := range res.Points {
		points = append(points, ivPointDTO{
			EtaV:   roundE(p.EtaV),
			IBV:    roundCurrent(p.IBV),
			ITafel: roundCurrent(p.ITafel),
		})
	}
	return ivResponse{
		EquilibriumPotentialV:  roundE(res.EquilibriumPotentialV),
		SlopeMVPerDecade:       roundSlope(res.SlopeMVPerDecade),
		TafelSlopeMVPerDecade:  roundSlope(res.TafelSlopeMVPerDecade),
		ExchangeCurrentDensity: roundCurrent(res.ExchangeCurrentDensity),
		AlphaAnodic:            roundSlope(res.AlphaAnodic),
		AlphaCathodic:          roundSlope(res.AlphaCathodic),
		Points:                 points,
	}
}
