package main

import (
	"log"
	"net/http"

	"validator-hub/service-cnpj/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)
	mux.HandleFunc("GET /cnpj/{cnpj}", handler.ValidatePath)

	log.Println("iniciando service-cnpj na porta 8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
