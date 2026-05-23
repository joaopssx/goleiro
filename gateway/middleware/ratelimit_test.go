package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRateLimitInternal(t *testing.T) {
	rl := NewRateLimit()
	
	for i := 0; i < 100; i++ {
		allowed, _ := rl.allow("127.0.0.1")
		if !allowed {
			t.Fatal("ip interno bloqueado")
		}
	}
}

func TestRateLimitConsumption(t *testing.T) {
	os.Setenv("RATE_LIMIT_RPM", "60")
	rl := NewRateLimit()

	for i := 0; i < 60; i++ {
		allowed, _ := rl.allow("1.2.3.4")
		if !allowed {
			t.Fatalf("bloqueado prematuramente na req %d", i+1)
		}
	}

	allowed, wait := rl.allow("1.2.3.4")
	if allowed {
		t.Fatal("deveria ter bloqueado a 61ª req")
	}
	if wait <= 0 {
		t.Fatal("wait duration deveria ser positivo")
	}
}

func TestRateLimitHandler(t *testing.T) {
	rl := NewRateLimit()
	handler := rl.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatal("esperava 200 por ser ip interno")
	}
}
