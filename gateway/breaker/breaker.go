package breaker

import (
	"net/http"
	"sync"
	"time"

	"validator-hub/gateway/logger"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu           sync.Mutex
	name         string
	state        State
	failures     int
	lastFailure  time.Time
	stateChanged time.Time
}

func New(name string) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		state:        StateClosed,
		stateChanged: time.Now(),
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures > 0 && now.Sub(cb.lastFailure) > 30*time.Second {
			cb.failures = 0
		}
		return true

	case StateOpen:
		if now.Sub(cb.stateChanged) > 20*time.Second {
			cb.setState(StateHalfOpen, now)
			return true
		}
		return false

	case StateHalfOpen:
		return false
	}

	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.setState(StateClosed, time.Now())
	}
	cb.failures = 0
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	cb.lastFailure = now

	if cb.state == StateHalfOpen {
		cb.setState(StateOpen, now)
		return
	}

	if cb.state == StateClosed {
		cb.failures++
		if cb.failures >= 5 {
			cb.setState(StateOpen, now)
		}
	}
}

func (cb *CircuitBreaker) setState(newState State, now time.Time) {
	old := cb.state
	cb.state = newState
	cb.stateChanged = now
	
	oldStr := "fechado"
	if old == StateOpen { oldStr = "aberto" }
	if old == StateHalfOpen { oldStr = "semi-aberto" }
	
	newStr := "fechado"
	if newState == StateOpen { newStr = "aberto" }
	if newState == StateHalfOpen { newStr = "semi-aberto" }

	logger.StateChange("circuit breaker alterado ("+cb.name+")", oldStr, newStr)
}

type BreakerTransport struct {
	Base    http.RoundTripper
	Breaker *CircuitBreaker
}

func (t *BreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.Breaker.Allow() {
		return nil, ErrCircuitOpen
	}

	resp, err := t.Base.RoundTrip(req)
	
	if err != nil || resp.StatusCode >= 500 {
		t.Breaker.RecordFailure()
	} else {
		t.Breaker.RecordSuccess()
	}

	return resp, err
}

type circuitError string
func (e circuitError) Error() string { return string(e) }
const ErrCircuitOpen = circuitError("circuito aberto")
