package main

import (
	"log"
	"net/http"

	"validator-hub/service-email/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)

	log.Println("iniciando service-email na porta 8084")
	if err := http.ListenAndServe(":8084", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
