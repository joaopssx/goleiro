package main

import (
	"log"
	"net/http"

	"validator-hub/service-cpf/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)
	mux.HandleFunc("GET /cpf/{cpf}", handler.ValidatePath)

	log.Println("iniciando service-cpf na porta 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
