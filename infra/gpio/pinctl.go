package gpio

import (
	"os/exec"
	"strings"
)

// GetAvailablePins Função para obter a lista de pinos disponíveis
func GetAvailablePins() ([]string, error) {
	output, err := exec.Command("pinctrl", "list").Output()
	if err != nil {
		return nil, err
	}
	// Convertendo a saída em uma lista de pinos
	pins := strings.Fields(string(output))
	return pins, nil
}

// GetPinState Função para obter o estado de um pino específico
func GetPinState(pin string) (string, error) {
	output, err := exec.Command("pinctrl", "get", pin).Output()
	if err != nil {
		return "", err
	}
	// Convertendo a saída em uma string
	state := strings.TrimSpace(string(output))
	return state, nil
}
