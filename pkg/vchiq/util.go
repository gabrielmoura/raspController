package vchiq

import (
	"regexp"
	"strings"
)

// clean removes specified substrings from the input string and trims the result.
func clean(str string, args ...string) string {
	for _, arg := range args {
		str = strings.ReplaceAll(str, arg, "")
	}
	return strings.TrimSpace(str)
}

// extractMemorySize extrai o número e a letra associada (M, G, ou K) de uma string.
// Retorna uma única string contendo o número e a letra.
func extractMemorySize(input string) string {
	// Expressão regular para encontrar o valor numérico seguido de M, G ou K
	re := regexp.MustCompile(`(\d+)([MGK])`)
	match := re.FindStringSubmatch(input)
	if len(match) >= 3 {
		// Concatenar o número e a letra em uma única string
		return match[1] + match[2]
	}
	// Se não encontrar correspondência, retorna uma string vazia
	return ""
}
