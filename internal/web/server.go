// Package web 提供 nernst-cell 的 HTTP 控制台：
// /api/nernst 与 /api/iv 两个求解接口、/api/examples 算例
// 列表，以及由 Go 同进程托管的静态页面。错误一律返回
// 带 error 字段的 JSON，页面与 API 都能读到后端失败原因。
package web

import (
	"io/fs"
	"net/http"
)

// Assets 是构造服务所需的静态资源与内置算例。
type Assets struct {
	WebFS    fs.FS
	Examples map[string][]byte
}

// Server 持有路由 mux 与静态资源。
type Server struct {
	assets Assets
	mux    *http.ServeMux
}

// NewServer 注册全部路由并返回可用的 http.Handler。
func NewServer(assets Assets) http.Handler {
	s := &Server{assets: assets, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/nernst", s.handleNernst)
	s.mux.HandleFunc("/api/iv", s.handleIV)
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.Handle("/", staticHandler(assets.WebFS))
	return withRequestLog(s.mux)
}

// ServeHTTP 转发到内部 mux。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
