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

	"validator-hub/service-email/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)

	srv := &http.Server{
		Addr:    ":8084",
		Handler: mux,
	}

	go func() {
		log.Println("iniciando service-email na porta 8084")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro ao iniciar servidor: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("encerrando serviço, aguardando requisições em andamento...")

	timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT_SECONDS")
	timeoutSec := 15
	if val, err := strconv.Atoi(timeoutStr); err == nil && val > 0 {
		timeoutSec = val
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Println("timeout de encerramento atingido, forçando saída")
		os.Exit(1)
	}

	log.Println("serviço encerrado com sucesso")
	os.Exit(0)
}
