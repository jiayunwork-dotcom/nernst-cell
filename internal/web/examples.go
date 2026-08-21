package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

// errExampleNotFound 在内置算例缺失时返回。
var errExampleNotFound = errors.New("example not found")

// examplesResponse 是 /api/examples 的响应体：
// 键为算例名，值为算例原始 JSON。
type examplesResponse struct {
	CuConc json.RawMessage `json:"cu-conc"`
}

// handleExamples 处理 GET /api/examples，返回内置算例供页面
// 一键加载。
func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	cuConc, ok := s.assets.Examples["cu-conc"]
	if !ok {
		writeError(w, http.StatusNotFound, errExampleNotFound)
		return
	}
	writeJSON(w, http.StatusOK, examplesResponse{
		CuConc: cuConc,
	})
}
