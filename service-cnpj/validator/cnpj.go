package validator

import (
	"regexp"
	"strconv"
)

var cleanRegex = regexp.MustCompile(`\D+`)

func ValidateCNPJ(input string) (bool, string, string) {
	cnpj := cleanRegex.ReplaceAllString(input, "")

	if len(cnpj) != 14 {
		return false, input, "cnpj inválido: tamanho incorreto"
	}

	if isRepeatedSequence(cnpj) {
		return false, formatCNPJ(cnpj), "cnpj inválido: sequência de dígitos repetidos"
	}

	if !verifyDigits(cnpj) {
		return false, formatCNPJ(cnpj), "cnpj inválido: dígitos verificadores incorretos"
	}

	return true, formatCNPJ(cnpj), "cnpj válido"
}

func isRepeatedSequence(cnpj string) bool {
	first := cnpj[0]
	for i := 1; i < 14; i++ {
		if cnpj[i] != first {
			return false
		}
	}
	return true
}

func verifyDigits(cnpj string) bool {
	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	d1 := calculateDigit(cnpj[:12], w1)
	if strconv.Itoa(d1) != string(cnpj[12]) {
		return false
	}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	d2 := calculateDigit(cnpj[:13], w2)
	if strconv.Itoa(d2) != string(cnpj[13]) {
		return false
	}
	return true
}

func calculateDigit(base string, weights []int) int {
	sum := 0
	for i, char := range base {
		num, _ := strconv.Atoi(string(char))
		sum += num * weights[i]
	}
	rem := sum % 11
	if rem < 2 {
		return 0
	}
	return 11 - rem
}

func formatCNPJ(cnpj string) string {
	if len(cnpj) == 14 {
		return cnpj[:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:]
	}
	return cnpj
}
