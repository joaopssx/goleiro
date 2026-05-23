package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"validator-hub/service-cep/metrics"
	"validator-hub/service-cep/validator"
)

var startTime = time.Now()

type ValidationResponse struct {
	Valido    bool   `json:"valido"`
	Formatado string `json:"formatado"`
	Mensagem  string `json:"mensagem"`
}

type ValidationRequest struct {
	CEP string `json:"cep"`
}

func ValidateQuery(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	processValidation(w, cep)
}

func ValidatePath(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	processValidation(w, cep)
}

func ValidateBody(w http.ResponseWriter, r *http.Request) {
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "corpo da requisição inválido"}`))
		return
	}
	processValidation(w, req.CEP)
}

func processValidation(w http.ResponseWriter, cep string) {
	start := time.Now()
	
	if cep == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "cep não informado"}`))
		
		metrics.Record(false, true, uint64(time.Since(start).Milliseconds()))
		return
	}

	valid, formatted, message := validator.ValidateCEP(cep)

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
		Nome:      "service-cep",
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
		"nome":      "service-cep",
		"endpoints": "GET /validate?cep={cep}, GET /cep/{cep}, POST /validate, GET /health, GET /metricas",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
