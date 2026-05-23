package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type serviceStatus struct {
	Nome   string `json:"nome"`
	Status string `json:"status"`
}

type healthResponse struct {
	Status   string          `json:"status"`
	Servicos []serviceStatus `json:"servicos"`
}

var services = []string{
	"http://service-cpf:8081/health",
	"http://service-cnpj:8082/health",
	"http://service-cep:8083/health",
	"http://service-email:8084/health",
}

func Check(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []serviceStatus{}

	client := http.Client{Timeout: 2 * time.Second}

	for _, url := range services {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			
			nome := "desconhecido"
			switch target {
			case "http://service-cpf:8081/health": nome = "service-cpf"
			case "http://service-cnpj:8082/health": nome = "service-cnpj"
			case "http://service-cep:8083/health": nome = "service-cep"
			case "http://service-email:8084/health": nome = "service-email"
			}

			status := "indisponível"
			resp, err := client.Get(target)
			if err == nil {
				if resp.StatusCode == http.StatusOK {
					status = "ok"
				}
				resp.Body.Close()
			}

			mu.Lock()
			results = append(results, serviceStatus{Nome: nome, Status: status})
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	overall := "ok"
	for _, res := range results {
		if res.Status != "ok" {
			overall = "parcial"
			break
		}
	}

	resp := healthResponse{
		Status:   overall,
		Servicos: results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
