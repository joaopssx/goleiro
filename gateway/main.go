package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxy(target string) http.HandlerFunc {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/validate"
		p.ServeHTTP(w, r)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "mensagem": "gateway e rotas ativas"}`))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /v1/cpf/validate", proxy("http://service-cpf:8081"))
	mux.HandleFunc("POST /v1/cnpj/validate", proxy("http://service-cnpj:8082"))
	mux.HandleFunc("POST /v1/cep/validate", proxy("http://service-cep:8083"))
	mux.HandleFunc("POST /v1/email/validate", proxy("http://service-email:8084"))

	log.Println("iniciando gateway na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("erro ao iniciar gateway: %v", err)
	}
}
