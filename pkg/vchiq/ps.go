package vchiq

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
)

// ProcessInfo contains information about a process.
type ProcessInfo struct {
	PID     string
	PPID    string
	Cmd     string
	CPU     string
	Mem     string
	Elapsed string
}

var (
	ErrProcessNotFound = errors.New("process not found")
	ErrGettingProcess  = errors.New("error getting process")
)

// ListProcesses list all processes.
func ListProcesses() ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-e", "-o", "pid,ppid,comm,%cpu,%mem,etime", "--no-headers")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting process:", err)
		return nil, ErrGettingProcess
	}

	lines := strings.Split(out.String(), "\n")

	var processes []ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)

		// If there are less than 6 fields, something is wrong, so we ignore it
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

	if len(processes) == 0 {
		return nil, ErrProcessNotFound
	}

	return processes, nil
}

// GetProcessByPid return a process by PID.
func GetProcessByPid(pid string) ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-o", "pid,ppid,cmd,%cpu,etime,%mem", "--pid", pid, "--no-headers")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Error getting process:", err)
		return nil, ErrGettingProcess
	}

	lines := strings.Split(out.String(), "\n")

	var processes []ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)

		// If there are less than 6 fields, something is wrong, so we ignore it
		if len(fields) < 6 {
			continue
		}

		processes = append(processes, ProcessInfo{
			PID:     fields[0],
			PPID:    fields[1],
			Cmd:     fields[2],
			CPU:     fields[3],
			Elapsed: fields[4],
			Mem:     fields[5],
		})
	}
	if len(processes) == 0 {
		return nil, ErrProcessNotFound
	}
	return processes, nil
}
