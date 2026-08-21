package web

import (
	"net/http"

	"nernst-cell/internal/nernst"
)

// handleNernst 处理 POST /api/nernst：E°、n、T、活度 → E 与斜率。
// 非法输入返回 400 与 error JSON。
func (s *Server) handleNernst(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req nernstRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	in, err := req.toInput()
	if err != nil {
		badRequest(w, err)
		return
	}
	res, err := nernst.Evaluate(in)
	if err != nil {
		badRequest(w, err)
		return
	}
	writeJSON(w, http.StatusOK, fromResult(res, in.Electrons))
}
