package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"validator-hub/service-cnpj/metrics"
	"validator-hub/service-cnpj/validator"
)

var startTime = time.Now()

type ValidationResponse struct {
	Valido    bool   `json:"valido"`
	Formatado string `json:"formatado"`
	Mensagem  string `json:"mensagem"`
}

type ValidationRequest struct {
	CNPJ string `json:"cnpj"`
}

func ValidateQuery(w http.ResponseWriter, r *http.Request) {
	cnpj := r.URL.Query().Get("cnpj")
	processValidation(w, cnpj)
}

func ValidatePath(w http.ResponseWriter, r *http.Request) {
	cnpj := r.PathValue("cnpj")
	processValidation(w, cnpj)
}

func ValidateBody(w http.ResponseWriter, r *http.Request) {
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "corpo da requisição inválido"}`))
		return
	}
	processValidation(w, req.CNPJ)
}

func processValidation(w http.ResponseWriter, cnpj string) {
	start := time.Now()
	
	if cnpj == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "cnpj não informado"}`))
		
		metrics.Record(false, true, uint64(time.Since(start).Milliseconds()))
		return
	}

	valid, formatted, message := validator.ValidateCNPJ(cnpj)

	resp := ValidationResponse{
		Valido:    valid,
		Formatado: formatted,
		Mensagem:  message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
	
	metrics.Record(valid, false, uint64(time.Since(start).Milliseconds()))
}

type HealthResponse struct {
	Nome      string `json:"nome"`
	Status    string `json:"status"`
	Versao    string `json:"versao"`
	Uptime    int64  `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Nome:      "service-cnpj",
		Status:    "ok",
		Versao:    "1.0.0",
		Uptime:    int64(time.Since(startTime).Seconds()),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func Root(w http.ResponseWriter, r *http.Request) {
	docs := map[string]string{
		"nome":      "service-cnpj",
		"endpoints": "GET /validate?cnpj={cnpj}, GET /cnpj/{cnpj}, POST /validate, GET /health, GET /metricas",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
