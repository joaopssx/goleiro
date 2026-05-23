package validator

import (
	"regexp"
)

var cleanRegex = regexp.MustCompile(`\D+`)

func ValidateCEP(input string) (bool, string, string) {
	cep := cleanRegex.ReplaceAllString(input, "")

	if len(cep) != 8 {
		return false, input, "cep inválido: tamanho incorreto"
	}

	return true, formatCEP(cep), "cep válido"
}

func formatCEP(cep string) string {
	if len(cep) == 8 {
		return cep[:5] + "-" + cep[5:]
	}
	return cep
}
