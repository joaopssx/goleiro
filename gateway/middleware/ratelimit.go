package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimit struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

func NewRateLimit() *RateLimit {
	rl := &RateLimit{
		requests: make(map[string][]time.Time),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimit) cleanup() {
	for {
		time.Sleep(2 * time.Minute)
		rl.mu.Lock()
		now := time.Now()
		for ip, times := range rl.requests {
			valid := []time.Time{}
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimit) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[ip]

	valid := []time.Time{}
	for _, t := range times {
		if now.Sub(t) < time.Minute {
			valid = append(valid, t)
		}
	}

	if len(valid) >= 60 {
		rl.requests[ip] = valid
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func (rl *RateLimit) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if idx := strings.LastIndex(ip, ":"); idx != -1 {
			ip = ip[:idx]
		}
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = strings.Split(xff, ",")[0]
		}

		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"erro": "limite de requisições excedido"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
