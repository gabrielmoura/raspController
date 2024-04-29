package vchiq

import (
	"bytes"
	"os/exec"
	"strings"
)

// ProcessInfo contém informações sobre um processo.
type ProcessInfo struct {
	PID     string
	PPID    string
	Cmd     string
	CPU     string
	Mem     string
	Elapsed string
}

// ListProcesses lista todos os processos e suas informações.
func ListProcesses() ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-e", "-o", "pid,ppid,comm,%cpu,%mem,etime")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")

	var processes []ProcessInfo
	for i, line := range lines {
		// Ignora o cabeçalho da tabela
		if i == 0 || len(line) == 0 {
			continue
		}

		fields := strings.Fields(line)

		// Se houver menos de 6 campos, algo está errado, então ignoramos
		if len(fields) < 6 {
			continue
		}

		processes = append(processes, ProcessInfo{
			PID:     fields[0],
			PPID:    fields[1],
			Cmd:     fields[2],
			CPU:     fields[3],
			Mem:     fields[4],
			Elapsed: fields[5],
		})
	}

	return processes, nil
}
