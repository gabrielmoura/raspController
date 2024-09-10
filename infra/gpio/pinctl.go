package gpio

import (
	"errors"
	"os/exec"
	"strings"
)

// GetAvailablePins get the list of available pins
func GetAvailablePins() ([]string, error) {
	if _, err := exec.LookPath("lsgpio"); err != nil {
		return nil, errors.New("pinctrl command not found")
	}

	output, err := exec.Command("pinctrl").Output()
	if err != nil {
		return nil, err
	}
	// Converting the output to a list of pins
	pins := strings.Fields(string(output))
	return pins, nil

}

// GetPinState get the state of a pin
func GetPinState(pin string) (string, error) {
	if _, err := exec.LookPath("lsgpio"); err != nil {
		return "", errors.New("pinctrl command not found")
	}
	output, err := exec.Command("pinctrl", "get", pin).Output()
	if err != nil {
		return "", err
	}
	// Converting the output to a string
	state := strings.TrimSpace(string(output))
	return state, nil
}
