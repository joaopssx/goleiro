package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ServiceDetails struct {
	Status     string `json:"status"`
	LatenciaMs int64  `json:"latencia_ms"`
}

type HealthResponse struct {
	Status       string                    `json:"status"`
	Servicos     map[string]ServiceDetails `json:"servicos"`
	VerificadoEm string                    `json:"verificado_em"`
}

type ServiceResponse struct {
	Nome      string `json:"nome"`
	Status    string `json:"status"`
	Versao    string `json:"versao"`
	Uptime    int64  `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

type DetailedServiceDetails struct {
	Status     string `json:"status"`
	LatenciaMs int64  `json:"latencia_ms"`
	Versao     string `json:"versao,omitempty"`
	Uptime     int64  `json:"uptime,omitempty"`
}

type DetailedHealthResponse struct {
	Status       string                            `json:"status"`
	Servicos     map[string]DetailedServiceDetails `json:"servicos"`
	VerificadoEm string                            `json:"verificado_em"`
}

var services = map[string]string{
	"cpf":   "http://service-cpf:8081/health",
	"cnpj":  "http://service-cnpj:8082/health",
	"cep":   "http://service-cep:8083/health",
	"email": "http://service-email:8084/health",
}

func getHealth(detailed bool) (int, any) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	basicServicos := make(map[string]ServiceDetails)
	detailedServicos := make(map[string]DetailedServiceDetails)

	client := http.Client{Timeout: 2 * time.Second}

	for key, url := range services {
		wg.Add(1)
		go func(k, u string) {
			defer wg.Done()

			start := time.Now()
			resp, err := client.Get(u)
			latencia := time.Since(start).Milliseconds()

			status := "inacessível"
			versao := ""
			var uptime int64 = 0

			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					status = "ok"
				} else if resp.StatusCode == http.StatusServiceUnavailable {
					status = "degradado"
				} else {
					status = "degradado"
				}

				var svcResp ServiceResponse
				if json.NewDecoder(resp.Body).Decode(&svcResp) == nil {
					versao = svcResp.Versao
					uptime = svcResp.Uptime
				}
			}

			mu.Lock()
			if detailed {
				detailedServicos[k] = DetailedServiceDetails{
					Status:     status,
					LatenciaMs: latencia,
					Versao:     versao,
					Uptime:     uptime,
				}
			} else {
				basicServicos[k] = ServiceDetails{
					Status:     status,
					LatenciaMs: latencia,
				}
			}
			mu.Unlock()
		}(key, url)
	}

	wg.Wait()

	overallStatus := "ok"
	statusCode := http.StatusOK

	if detailed {
		for _, s := range detailedServicos {
			if s.Status == "inacessível" {
				overallStatus = "inacessível"
				statusCode = http.StatusServiceUnavailable
			} else if s.Status == "degradado" && overallStatus != "inacessível" {
				overallStatus = "degradado"
				statusCode = 207
			}
		}
		return statusCode, DetailedHealthResponse{
			Status:       overallStatus,
			Servicos:     detailedServicos,
			VerificadoEm: time.Now().UTC().Format(time.RFC3339),
		}
	} else {
		for _, s := range basicServicos {
			if s.Status == "inacessível" {
				overallStatus = "inacessível"
				statusCode = http.StatusServiceUnavailable
			} else if s.Status == "degradado" && overallStatus != "inacessível" {
				overallStatus = "degradado"
				statusCode = 207
			}
		}
		return statusCode, HealthResponse{
			Status:       overallStatus,
			Servicos:     basicServicos,
			VerificadoEm: time.Now().UTC().Format(time.RFC3339),
		}
	}
}

func Check(w http.ResponseWriter, r *http.Request) {
	code, payload := getHealth(false)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func CheckDetailed(w http.ResponseWriter, r *http.Request) {
	code, payload := getHealth(true)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
