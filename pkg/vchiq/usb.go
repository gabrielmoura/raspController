package vchiq

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

type USBDevice struct {
	Bus     string
	Device  string
	ID      string
	Vendor  string
	Product string
}

var (
	ErrExecutingLsusb = errors.New("error executing lsusb")
)

func GetUsbList() ([]USBDevice, error) {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return nil, ErrExecutingLsusb
	}
	// Converte a saída para string
	output := string(out)

	// Define uma expressão regular para extrair as informações
	re := regexp.MustCompile(`Bus (\d+) Device (\d+): ID (\w+):(\w+) (.+)`)

	// Cria um slice para armazenar os dispositivos USB
	var devices []USBDevice

	// Percorre cada linha da saída
	for _, line := range strings.Split(output, "\n") {
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			devices = append(devices, USBDevice{
				Bus:     matches[1],
				Device:  matches[2],
				ID:      matches[3] + ":" + matches[4],
				Vendor:  matches[3],
				Product: matches[4] + " " + matches[5],
			})
		}
	}
	return devices, nil
}
