package handler

import (
	"encoding/json"
	"net/http"

	"validator-hub/service-cpf/validator"
)

type ValidationResponse struct {
	Valido    bool   `json:"valido"`
	Formatado string `json:"formatado"`
	Mensagem  string `json:"mensagem"`
}

type ValidationRequest struct {
	CPF string `json:"cpf"`
}

func ValidateQuery(w http.ResponseWriter, r *http.Request) {
	cpf := r.URL.Query().Get("cpf")
	processValidation(w, cpf)
}

func ValidatePath(w http.ResponseWriter, r *http.Request) {
	cpf := r.PathValue("cpf")
	processValidation(w, cpf)
}

func ValidateBody(w http.ResponseWriter, r *http.Request) {
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "corpo da requisição inválido"}`))
		return
	}
	processValidation(w, req.CPF)
}

func processValidation(w http.ResponseWriter, cpf string) {
	if cpf == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "cpf não informado"}`))
		return
	}

	valid, formatted, message := validator.ValidateCPF(cpf)

	resp := ValidationResponse{
		Valido:    valid,
		Formatado: formatted,
		Mensagem:  message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func Root(w http.ResponseWriter, r *http.Request) {
	docs := map[string]string{
		"nome":      "service-cpf",
		"endpoints": "GET /validate?cpf={cpf}, GET /cpf/{cpf}, POST /validate, GET /health",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
