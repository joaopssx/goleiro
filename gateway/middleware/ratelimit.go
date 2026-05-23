package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	mu     sync.Mutex
	tokens float64
	last   time.Time
}

type RateLimit struct {
	buckets sync.Map
	rate    float64
	limit   float64
}

func isInternalIP(ipStr string) bool {
	if ipStr == "127.0.0.1" || ipStr == "::1" {
		return true
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	if ip.IsPrivate() || ip.IsLoopback() {
		return true
	}
	return false
}

func NewRateLimit() *RateLimit {
	rpm := 60.0
	if val, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT_RPM"), 64); err == nil && val > 0 {
		rpm = val
	}
	rl := &RateLimit{
		rate:  rpm / 60.0,
		limit: rpm,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimit) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		now := time.Now()
		rl.buckets.Range(func(key, value any) bool {
			b := value.(*bucket)
			b.mu.Lock()
			if now.Sub(b.last) > 5*time.Minute {
				rl.buckets.Delete(key)
			}
			b.mu.Unlock()
			return true
		})
	}
}

func (rl *RateLimit) allow(ip string) (bool, time.Duration) {
	if isInternalIP(ip) {
		return true, 0
	}

	val, _ := rl.buckets.LoadOrStore(ip, &bucket{
		tokens: rl.limit,
		last:   time.Now(),
	})

	b := val.(*bucket)
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.limit {
		b.tokens = rl.limit
	}
	b.last = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return true, 0
	}

	waitSecs := (1 - b.tokens) / rl.rate
	return false, time.Duration(waitSecs * float64(time.Second))
}

type RateLimitResponse struct {
	Erro           string `json:"erro"`
	Limite         int    `json:"limite"`
	Janela         string `json:"janela"`
	TenteNovamente string `json:"tente_novamente_em"`
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
		ip = strings.TrimSpace(ip)

		allowed, wait := rl.allow(ip)
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())+1))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			
			resp := RateLimitResponse{
				Erro:           "limite de requisições excedido",
				Limite:         int(rl.limit),
				Janela:         "1 minuto",
				TenteNovamente: time.Now().Add(wait).UTC().Format(time.RFC3339),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
