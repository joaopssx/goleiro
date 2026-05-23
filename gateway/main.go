package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"validator-hub/gateway/router"
)

func main() {
	log.SetFlags(0)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal(`{"erro": "variável de ambiente API_KEY não configurada"}`)
	}

	mux := router.New(apiKey)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println(`{"mensagem": "iniciando gateway na porta 8080"}`)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf(`{"erro": "falha ao iniciar gateway", "detalhe": "%v"}`, err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println(`{"mensagem": "encerrando serviço, aguardando requisições em andamento..."}`)

	timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT_SECONDS")
	timeoutSec := 15
	if val, err := strconv.Atoi(timeoutStr); err == nil && val > 0 {
		timeoutSec = val
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Println(`{"erro": "timeout de encerramento atingido, forçando saída"}`)
		os.Exit(1)
	}

	log.Println(`{"mensagem": "serviço encerrado com sucesso"}`)
	os.Exit(0)
}
