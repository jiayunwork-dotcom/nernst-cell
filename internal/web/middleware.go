package web

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder 记录响应状态码，供请求日志使用。
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader 覆盖底层写入并记录状态码。
func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// withRequestLog 打印每个请求的方法、路径、状态码与耗时。
func withRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
	})
}
