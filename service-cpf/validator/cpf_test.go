package validator

import "testing"

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		valid     bool
		formatted string
		msg       string
	}{
		{"cpf válido com máscara", "123.456.789-09", true, "123.456.789-09", "cpf válido"},
		{"cpf válido sem máscara", "12345678909", true, "123.456.789-09", "cpf válido"},
		{"cpf inválido tamanho incorreto", "123", false, "123", "cpf inválido: tamanho incorreto"},
		{"cpf inválido sequência repetida", "111.111.111-11", false, "111.111.111-11", "cpf inválido: sequência de dígitos repetidos"},
		{"cpf inválido dígitos incorretos", "123.456.789-00", false, "123.456.789-00", "cpf inválido: dígitos verificadores incorretos"},
		{"entrada vazia", "", false, "", "cpf inválido: tamanho incorreto"},
		{"entrada com letras", "abc.def.ghi-jk", false, "abc.def.ghi-jk", "cpf inválido: tamanho incorreto"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, formatted, msg := ValidateCPF(tt.input)
			if valid != tt.valid {
				t.Errorf("esperava valido=%v, obteve %v", tt.valid, valid)
			}
			if formatted != tt.formatted {
				t.Errorf("esperava formatado=%q, obteve %q", tt.formatted, formatted)
			}
			if msg != tt.msg {
				t.Errorf("esperava msg=%q, obteve %q", tt.msg, msg)
			}
		})
	}
}
