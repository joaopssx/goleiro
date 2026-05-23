package validator

import (
	"regexp"
	"strconv"
)

var cleanRegex = regexp.MustCompile(`\D+`)

func ValidateCPF(input string) (bool, string, string) {
	cpf := cleanRegex.ReplaceAllString(input, "")

	if len(cpf) != 11 {
		return false, input, "cpf inválido: tamanho incorreto"
	}

	if isRepeatedSequence(cpf) {
		return false, formatCPF(cpf), "cpf inválido: sequência de dígitos repetidos"
	}

	if !verifyDigits(cpf) {
		return false, formatCPF(cpf), "cpf inválido: dígitos verificadores incorretos"
	}

	return true, formatCPF(cpf), "cpf válido"
}

func isRepeatedSequence(cpf string) bool {
	first := cpf[0]
	for i := 1; i < 11; i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

func verifyDigits(cpf string) bool {
	d1 := calculateDigit(cpf[:9])
	if strconv.Itoa(d1) != string(cpf[9]) {
		return false
	}
	d2 := calculateDigit(cpf[:10])
	if strconv.Itoa(d2) != string(cpf[10]) {
		return false
	}
	return true
}

func calculateDigit(base string) int {
	sum := 0
	weight := len(base) + 1
	for _, char := range base {
		num, _ := strconv.Atoi(string(char))
		sum += num * weight
		weight--
	}
	rem := sum % 11
	if rem < 2 {
		return 0
	}
	return 11 - rem
}

func formatCPF(cpf string) string {
	if len(cpf) == 11 {
		return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
	}
	return cpf
}
