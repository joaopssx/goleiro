package main

import (
	"log"
	"net/http"
	"os"

	"validator-hub/gateway/router"
)

func main() {
	log.SetFlags(0) // disable standard timestamp since we use structured json logging

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal(`{"erro": "variável de ambiente API_KEY não configurada"}`)
	}

	mux := router.New(apiKey)

	log.Println(`{"mensagem": "iniciando gateway na porta 8080"}`)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf(`{"erro": "falha ao iniciar gateway", "detalhe": "%v"}`, err)
	}
}
