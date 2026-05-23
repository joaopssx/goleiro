package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type logEntry struct {
	Metodo string `json:"metodo"`
	Path   string `json:"path"`
	Ip     string `json:"ip"`
	Tempo  string `json:"tempo"`
	Status int    `json:"status"`
}

type responseWriterObserver struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterObserver) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ip := r.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = strings.Split(xff, ",")[0]
		}
		
		rw := &responseWriterObserver{ResponseWriter: w, status: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		entry := logEntry{
			Metodo: r.Method,
			Path:   r.URL.Path,
			Ip:     ip,
			Tempo:  time.Since(start).String(),
			Status: rw.status,
		}
		
		b, _ := json.Marshal(entry)
		log.Println(string(b))
	})
}
