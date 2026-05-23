package router

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"validator-hub/gateway/breaker"
	"validator-hub/gateway/health"
	"validator-hub/gateway/middleware"
)

func New(apiKey string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.Check)
	mux.HandleFunc("GET /health/detalhado", health.CheckDetailed)

	proxyCpf := createProxy("http://service-cpf:8081", "cpf")
	proxyCnpj := createProxy("http://service-cnpj:8082", "cnpj")
	proxyCep := createProxy("http://service-cep:8083", "cep")
	proxyEmail := createProxy("http://service-email:8084", "email")

	mux.Handle("/cpf/", http.StripPrefix("/cpf", proxyCpf))
	mux.Handle("/cnpj/", http.StripPrefix("/cnpj", proxyCnpj))
	mux.Handle("/cep/", http.StripPrefix("/cep", proxyCep))
	mux.Handle("/email/", http.StripPrefix("/email", proxyEmail))

	authMw := middleware.NewAuth(apiKey)
	rateMw := middleware.NewRateLimit()

	return middleware.Logger(rateMw.Handler(authMw.Handler(mux)))
}

type ErrorResponse struct {
	Erro    string `json:"erro"`
	Servico string `json:"servico"`
}

func createProxy(target string, serviceName string) http.Handler {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)

	cb := breaker.New(serviceName)

	p.Transport = &breaker.BreakerTransport{
		Base:    http.DefaultTransport,
		Breaker: cb,
	}

	p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		resp := ErrorResponse{
			Erro:    "serviço temporariamente indisponível",
			Servico: serviceName,
		}
		json.NewEncoder(w).Encode(resp)
	}

	return p
}
