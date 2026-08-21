package web

import (
	"net/http"

	"nernst-cell/internal/polar"
)

// handleIV 处理 POST /api/iv：Nernst 输入 + i0、α、η 网格
// → i(η) 点列与 Tafel 对照。非法输入返回 400 与 error JSON。
func (s *Server) handleIV(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req ivRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	in, err := req.toIVInput()
	if err != nil {
		badRequest(w, err)
		return
	}
	res, err := polar.BuildIVCurve(in)
	if err != nil {
		badRequest(w, err)
		return
	}
	writeJSON(w, http.StatusOK, fromIVResult(res))
}
