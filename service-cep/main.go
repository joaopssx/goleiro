package main

import (
	"log"
	"net/http"

	"validator-hub/service-cep/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)
	mux.HandleFunc("GET /cep/{cep}", handler.ValidatePath)

	log.Println("iniciando service-cep na porta 8083")
	if err := http.ListenAndServe(":8083", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
