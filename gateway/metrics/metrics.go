package metrics

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type ServiceMetricas struct {
	TotalRequisicoes uint64 `json:"total_requisicoes"`
	TotalValidas     uint64 `json:"total_validas"`
	TotalInvalidas   uint64 `json:"total_invalidas"`
	TotalErros       uint64 `json:"total_erros"`
	TempoMedioMs     uint64 `json:"tempo_medio_ms"`
}

type Consolidado struct {
	TotalRequisicoes uint64                     `json:"total_requisicoes"`
	TotalValidas     uint64                     `json:"total_validas"`
	TotalInvalidas   uint64                     `json:"total_invalidas"`
	TotalErros       uint64                     `json:"total_erros"`
	Servicos         map[string]ServiceMetricas `json:"servicos"`
}

var (
	GatewayRequisicoes uint64
	GatewayTempoMs     uint64
	GatewayErros       uint64
)

func Handler(w http.ResponseWriter, r *http.Request) {
	services := map[string]string{
		"cpf":   "http://service-cpf:8081/metricas",
		"cnpj":  "http://service-cnpj:8082/metricas",
		"cep":   "http://service-cep:8083/metricas",
		"email": "http://service-email:8084/metricas",
	}

	client := http.Client{Timeout: 2 * time.Second}
	var wg sync.WaitGroup
	var mu sync.Mutex

	respMap := make(map[string]ServiceMetricas)

	for name, url := range services {
		wg.Add(1)
		go func(n, u string) {
			defer wg.Done()
			resp, err := client.Get(u)
			if err == nil {
				defer resp.Body.Close()
				var sm ServiceMetricas
				if json.NewDecoder(resp.Body).Decode(&sm) == nil {
					mu.Lock()
					respMap[n] = sm
					mu.Unlock()
				}
			}
		}(name, url)
	}

	wg.Wait()

	var totalReq, totalVal, totalInv, totalErr uint64
	for _, sm := range respMap {
		totalReq += sm.TotalRequisicoes
		totalVal += sm.TotalValidas
		totalInv += sm.TotalInvalidas
		totalErr += sm.TotalErros
	}

	gwReq := atomic.LoadUint64(&GatewayRequisicoes)
	gwErr := atomic.LoadUint64(&GatewayErros)
	gwTempo := atomic.LoadUint64(&GatewayTempoMs)
	
	gwMedio := uint64(0)
	if gwReq > 0 {
		gwMedio = gwTempo / gwReq
	}

	respMap["gateway"] = ServiceMetricas{
		TotalRequisicoes: gwReq,
		TotalErros:       gwErr,
		TempoMedioMs:     gwMedio,
	}
	
	totalReq += gwReq
	totalErr += gwErr

	consolidado := Consolidado{
		TotalRequisicoes: totalReq,
		TotalValidas:     totalVal,
		TotalInvalidas:   totalInv,
		TotalErros:       totalErr,
		Servicos:         respMap,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consolidado)
}

func Record(duracaoMs uint64, err bool) {
	atomic.AddUint64(&GatewayRequisicoes, 1)
	atomic.AddUint64(&GatewayTempoMs, duracaoMs)
	if err {
		atomic.AddUint64(&GatewayErros, 1)
	}
}
