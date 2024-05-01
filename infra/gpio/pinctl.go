package gpio

import (
	"errors"
	"os/exec"
	"strings"
)

// GetAvailablePins Função para obter a lista de pinos disponíveis
func GetAvailablePins() ([]string, error) {
	if _, err := exec.LookPath("lsgpio"); err != nil {
		return nil, errors.New("pinctrl command not found")
	}

	output, err := exec.Command("pinctrl").Output()
	if err != nil {
		return nil, err
	}
	// Convertendo a saída em uma lista de pinos
	pins := strings.Fields(string(output))
	return pins, nil

}

// GetPinState Função para obter o estado de um pino específico
func GetPinState(pin string) (string, error) {
	if _, err := exec.LookPath("lsgpio"); err != nil {
		return "", errors.New("pinctrl command not found")
	}
	output, err := exec.Command("pinctrl", "get", pin).Output()
	if err != nil {
		return "", err
	}
	// Convertendo a saída em uma string
	state := strings.TrimSpace(string(output))
	return state, nil
}
