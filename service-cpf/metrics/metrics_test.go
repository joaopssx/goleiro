package metrics

import (
	"sync"
	"testing"
)

func TestMetricsConcurrency(t *testing.T) {
	// reset counters
	TotalRequisicoes = 0
	TotalValidas = 0
	TotalInvalidas = 0
	TotalErros = 0
	TempoTotalMs = 0

	var wg sync.WaitGroup
	workers := 100
	requestsPerWorker := 1000

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < requestsPerWorker; j++ {
				if j%10 == 0 {
					Record(false, true, 10) // Erro
				} else if j%2 == 0 {
					Record(true, false, 5) // Valido
				} else {
					Record(false, false, 8) // Invalido
				}
			}
		}(i)
	}

	wg.Wait()

	m := GetMetricas()

	expectedTotal := uint64(workers * requestsPerWorker)
	if m.TotalRequisicoes != expectedTotal {
		t.Fatalf("esperado %d requisições, obtido %d", expectedTotal, m.TotalRequisicoes)
	}

	// Calculate expected values based on modulo logic
	// per worker (1000 requests):
	// j%10 == 0 -> 100 errors
	// j%2 == 0 (and not j%10==0) -> 400 valids
	// else -> 500 invalids
	
	expectedErros := uint64(workers * 100)
	expectedValidas := uint64(workers * 400)
	expectedInvalidas := uint64(workers * 500)

	if m.TotalErros != expectedErros {
		t.Fatalf("esperado %d erros, obtido %d", expectedErros, m.TotalErros)
	}
	if m.TotalValidas != expectedValidas {
		t.Fatalf("esperado %d válidas, obtido %d", expectedValidas, m.TotalValidas)
	}
	if m.TotalInvalidas != expectedInvalidas {
		t.Fatalf("esperado %d inválidas, obtido %d", expectedInvalidas, m.TotalInvalidas)
	}
}
