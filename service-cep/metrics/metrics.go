package metrics

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

var (
	TotalRequisicoes   uint64
	TotalValidas       uint64
	TotalInvalidas     uint64
	TotalErros         uint64
	TempoTotalMs       uint64
)

type MetricasResponse struct {
	TotalRequisicoes uint64 `json:"total_requisicoes"`
	TotalValidas     uint64 `json:"total_validas"`
	TotalInvalidas   uint64 `json:"total_invalidas"`
	TotalErros       uint64 `json:"total_erros"`
	TempoMedioMs     uint64 `json:"tempo_medio_ms"`
}

func GetMetricas() MetricasResponse {
	reqs := atomic.LoadUint64(&TotalRequisicoes)
	tempo := atomic.LoadUint64(&TempoTotalMs)
	
	medio := uint64(0)
	if reqs > 0 {
		medio = tempo / reqs
	}

	return MetricasResponse{
		TotalRequisicoes: reqs,
		TotalValidas:     atomic.LoadUint64(&TotalValidas),
		TotalInvalidas:   atomic.LoadUint64(&TotalInvalidas),
		TotalErros:       atomic.LoadUint64(&TotalErros),
		TempoMedioMs:     medio,
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := GetMetricas()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

func Record(valida bool, isErro bool, duracaoMs uint64) {
	atomic.AddUint64(&TotalRequisicoes, 1)
	atomic.AddUint64(&TempoTotalMs, duracaoMs)
	if isErro {
		atomic.AddUint64(&TotalErros, 1)
	} else if valida {
		atomic.AddUint64(&TotalValidas, 1)
	} else {
		atomic.AddUint64(&TotalInvalidas, 1)
	}
}
