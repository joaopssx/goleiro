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

	"validator-hub/service-cnpj/handler"
	"validator-hub/service-cnpj/logger"
	"validator-hub/service-cnpj/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		duracao := time.Since(start).Milliseconds()
		logger.Request("requisição processada", r.Method, r.URL.Path, r.RemoteAddr, rw.status, duracao)
	})
}

func main() {
	logger.Init("service-cnpj")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handler.Root)
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /metricas", metrics.Handler)
	
	mux.HandleFunc("GET /validate", handler.ValidateQuery)
	mux.HandleFunc("POST /validate", handler.ValidateBody)
	mux.HandleFunc("GET /cnpj/{cnpj}", handler.ValidatePath)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: loggingMiddleware(mux),
	}

	go func() {
		logger.Info("iniciando service-cnpj na porta 8082")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("erro ao iniciar servidor", "main.go", 59, err)
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
