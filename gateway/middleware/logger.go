package middleware

import (
	"net/http"
	"strings"
	"time"

	"validator-hub/gateway/logger"
	"validator-hub/gateway/metrics"
)

type responseWriterObserver struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterObserver) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriterObserver{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duracao := time.Since(start).Milliseconds()

		ip := r.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = strings.Split(xff, ",")[0]
		}
		ip = strings.TrimSpace(ip)

		logger.Request("requisição processada", r.Method, r.URL.Path, ip, rw.status, duracao)
		
		metrics.Record(uint64(duracao), rw.status >= 500)
	})
}
