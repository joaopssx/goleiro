package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"validator-hub/gateway/health"
	"validator-hub/gateway/middleware"
)

func New(apiKey string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.Check)
	mux.HandleFunc("GET /health/detalhado", health.CheckDetailed)

	proxyCpf := createProxy("http://service-cpf:8081")
	proxyCnpj := createProxy("http://service-cnpj:8082")
	proxyCep := createProxy("http://service-cep:8083")
	proxyEmail := createProxy("http://service-email:8084")

	mux.Handle("/cpf/", http.StripPrefix("/cpf", proxyCpf))
	mux.Handle("/cnpj/", http.StripPrefix("/cnpj", proxyCnpj))
	mux.Handle("/cep/", http.StripPrefix("/cep", proxyCep))
	mux.Handle("/email/", http.StripPrefix("/email", proxyEmail))

	authMw := middleware.NewAuth(apiKey)
	rateMw := middleware.NewRateLimit()

	return middleware.Logger(rateMw.Handler(authMw.Handler(mux)))
}

func createProxy(target string) http.Handler {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}
