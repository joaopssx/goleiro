package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"validator-hub/gateway/logger"
	"validator-hub/gateway/router"
)

func main() {
	logger.Init("gateway")

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logger.Error("variável de ambiente API_KEY não configurada", "main.go", 21, nil)
		os.Exit(1)
	}

	mux := router.New(apiKey)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Info("iniciando gateway na porta 8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("falha ao iniciar gateway", "main.go", 35, err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.StateChange("encerrando serviço, aguardando requisições em andamento...", "rodando", "encerrando")

	timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT_SECONDS")
	timeoutSec := 15
	if val, err := strconv.Atoi(timeoutStr); err == nil && val > 0 {
		timeoutSec = val
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.StateChange("timeout de encerramento atingido, forçando saída", "encerrando", "forçado")
		os.Exit(1)
	}

	logger.StateChange("serviço encerrado com sucesso", "encerrando", "encerrado")
	os.Exit(0)
}
