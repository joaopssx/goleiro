package handler

import (
	"encoding/json"
	"net/http"

	"validator-hub/service-cep/validator"
)

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
	if cep == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"valido": false, "formatado": "", "mensagem": "cep não informado"}`))
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
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func Root(w http.ResponseWriter, r *http.Request) {
	docs := map[string]string{
		"nome":      "service-cep",
		"endpoints": "GET /validate?cep={cep}, GET /cep/{cep}, POST /validate, GET /health",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
