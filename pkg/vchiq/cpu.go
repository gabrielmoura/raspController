package vchiq

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// Field representa um campo no JSON retornado pelo comando `lscpu -J`
type Field struct {
	Field    string  `json:"field"`
	Data     string  `json:"data"`
	Children []Field `json:"children,omitempty"`
}

// LscpuData estrutura para o JSON retornado pelo comando `lscpu -J`
type LscpuData struct {
	Lscpu []Field `json:"lscpu"`
}

var (
	ErrExecutingLscpu = errors.New("error executing lscpu")
	ErrParsingLscpu   = errors.New("error parsing lscpu output")
)

// GetCpus executa o comando `lscpu -J` e retorna os dados dos CPUs em formato []Field
func GetCpus() ([]Field, error) {
	cmd := exec.Command("lscpu", "-J")
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrExecutingLscpu
	}

	// Faz o parse da saída JSON
	var lscpuData LscpuData
	err = json.Unmarshal(output, &lscpuData)
	if err != nil {
		return nil, ErrParsingLscpu
	}

	return lscpuData.Lscpu, nil
}

// GetCPUCurrFreq retorna a frequência atual do CPU em MHz
func GetCPUCurrFreq() (float64, error) {
	out, err := exec.Command("vcgencmd", "measure_clock", "arm").Output()
	if err != nil {
		return 0, err
	}
	freqStr := strings.Split(string(out), "=")[1]
	freq, err := strconv.ParseFloat(strings.TrimSpace(freqStr), 64)
	if err != nil {
		return 0, err
	}
	return freq / 1000000, nil
}
