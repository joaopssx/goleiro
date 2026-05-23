package breaker

import (
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	cb := New("teste")

	for i := 0; i < 5; i++ {
		if !cb.Allow() {
			t.Fatal("deveria permitir")
		}
		cb.RecordFailure()
	}

	if cb.Allow() {
		t.Fatal("deveria bloquear após 5 falhas")
	}

	cb.mu.Lock()
	cb.stateChanged = time.Now().Add(-21 * time.Second)
	cb.mu.Unlock()

	if !cb.Allow() {
		t.Fatal("deveria permitir em semi-aberto")
	}

	if cb.Allow() {
		t.Fatal("não deveria permitir a segunda em semi-aberto")
	}

	cb.RecordFailure()

	if cb.Allow() {
		t.Fatal("deveria bloquear novamente após falha em semi-aberto")
	}

	cb.mu.Lock()
	cb.stateChanged = time.Now().Add(-21 * time.Second)
	cb.mu.Unlock()

	cb.Allow()
	
	cb.RecordSuccess()

	if !cb.Allow() {
		t.Fatal("deveria estar fechado após sucesso")
	}
}
