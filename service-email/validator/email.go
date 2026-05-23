package validator

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(input string) (bool, string, string) {
	if !emailRegex.MatchString(input) {
		return false, input, "email inválido: formato incorreto"
	}
	return true, input, "email válido"
}
