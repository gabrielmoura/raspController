// Package vchiq provides functions to retrieve various system information
// using vcgencmd and other system interfaces on Raspberry Pi.
package vchiq

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// Flags for various throttling and voltage events
	UnderVoltage          int64 = 1 << 0
	FreqCap                     = 1 << 1
	Throttling                  = 1 << 2
	SoftTempLimitActive         = 1 << 3
	UnderVoltageOccurred        = 1 << 16
	FreqCapOccurred             = 1 << 17
	Throttled                   = 1 << 18
	SoftTempLimitOccurred       = 1 << 19

	// Consulte https://www.raspberrypi.org/documentation/hardware/raspberrypi/revision-codes/README.md
	RpiZero    = "Raspberry Pi Zero"
	RpiZeroW   = "Raspberry Pi Zero W"
	Rpi3APlus  = "Raspberry Pi 3 Model A+"
	Rpi2B      = "Raspberry Pi 2 Model B"
	Rpi3B      = "Raspberry Pi 3 Model B"
	Rpi3BPlus  = "Raspberry Pi 3 Model B+"
	RpiCM3     = "Raspberry Pi Compute Module 3"
	RpiCM3Plus = "Raspberry Pi Compute Module 3+"
	Rpi4B      = "Raspberry Pi 4 Model B"
	Rpi400     = "Raspberry Pi 400"
	RpiCM4     = "Raspberry Pi Compute Module 4"
	RpiZero2W  = "Raspberry Pi Zero 2 W"
	Rpi5       = "Raspberry Pi 5"
)

// Mapeamento das revisões de CPU para os nomes dos dispositivos.
var deviceRevisions = map[string]string{
	"900021": RpiZero,
	"900032": RpiZeroW,
	"900092": RpiZero, "920092": RpiZero,
	"900093": RpiZero, "920093": RpiZero,
	"9000c1": RpiZeroW,
	"9020e0": Rpi3APlus, "9020e1": Rpi3APlus,
	"a01040": Rpi2B, "a01041": Rpi2B, "a21041": Rpi2B,
	"a02082": Rpi3B, "a22082": Rpi3B, "a32082": Rpi3B, "a52082": Rpi3B,
	"a020a0": RpiCM3, "a220a0": RpiCM3,
	"a020d3": Rpi3BPlus, "a22083": Rpi3BPlus, "a020d4": Rpi3BPlus,
	"a02042": Rpi2B + " (with BCM2837)", "a22042": Rpi2B + " (with BCM2837)",
	"a02100": RpiCM3Plus,
	"a03111": Rpi4B, "b03111": Rpi4B, "c03111": Rpi4B,
	"b03112": Rpi4B, "c03112": Rpi4B,
	"b03114": Rpi4B, "c03114": Rpi4B,
	"b03115": Rpi4B, "c03115": Rpi4B,
	"c03130": Rpi400,
	"a03140": RpiCM4, "b03140": RpiCM4, "c03140": RpiCM4, "d03140": RpiCM4,
	"902120": RpiZero2W,
	"c04170": Rpi5, "d04170": Rpi5,
}

// Mapeamento dos nomes dos dispositivos para a potência mínima do fornecimento de energia.
var minimalPowerSupply = map[string]float64{
	RpiZero:    1.2,
	RpiZeroW:   1.2,
	Rpi3APlus:  2.5,
	Rpi2B:      1.8,
	Rpi3B:      2.5,
	Rpi3BPlus:  2.5,
	RpiCM3:     2.5,
	RpiCM3Plus: 2.5,
	Rpi4B:      3.0,
	Rpi400:     3.0,
	RpiZero2W:  1.2,
	Rpi5:       3.0,
}

// GetThrottled returns the throttled status as an integer.
func GetThrottled() (int64, error) {
	rawThrottled, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return 0, fmt.Errorf("couldn't run vcgencmd: %w", err)
	}
	throttled, err := strconv.ParseInt(strings.TrimSpace(string(rawThrottled[12:])), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse throttled output: %w", err)
	}
	return throttled, nil
}

// GetThrottledInfo returns the throttled status as a string.
func GetThrottledInfo() (string, error) {
	throttled, err := GetThrottled()
	if err != nil {
		return "", err
	}

	var events []string

	if throttled&UnderVoltage != 0 {
		events = append(events, "Under-voltage detected")
	}
	if throttled&FreqCap != 0 {
		events = append(events, "Frequency capped")
	}
	if throttled&Throttling != 0 {
		events = append(events, "Throttling")
	}
	if throttled&SoftTempLimitActive != 0 {
		events = append(events, "Soft temperature limit active")
	}
	if throttled&UnderVoltageOccurred != 0 {
		events = append(events, "Under-voltage occurred")
	}
	if throttled&FreqCapOccurred != 0 {
		events = append(events, "Frequency cap occurred")
	}
	if throttled&Throttled != 0 {
		events = append(events, "Throttling occurred")
	}
	if throttled&SoftTempLimitOccurred != 0 {
		events = append(events, "Soft temperature limit occurred")
	}

	if len(events) == 0 {
		return "No throttling or voltage issues detected", nil
	}

	return strings.Join(events, "; "), nil
}

// GetGPUTemp returns the GPU temperature as a string.
func GetGPUTemp() (string, error) {
	temp, err := exec.Command("vcgencmd", "measure_temp").Output()
	if err != nil {
		return "", errors.New("couldn't run vcgencmd")
	}
	return clean(string(temp), "temp=", "'C"), nil
}

// GetCoreVolt returns the CPU voltage as a string.
func GetCoreVolt() (string, error) {
	volt, err := exec.Command("vcgencmd", "measure_volts").Output()
	if err != nil {
		return "", errors.New("couldn't run vcgencmd")
	}
	return clean(string(volt), "volt=", "V"), nil
}

// GetMem returns usage information for ARM and GPU memory as strings.
func GetMem() (string, string, error) {
	usageMem, err := exec.Command("vcgencmd", "get_mem", "arm").Output()
	if err != nil {
		return "", "", errors.New("couldn't run vcgencmd for arm memory")
	}
	gpuMem, err := exec.Command("vcgencmd", "get_mem", "gpu").Output()
	if err != nil {
		return "", "", errors.New("couldn't run vcgencmd for gpu memory")
	}
	return extractMemorySize(string(usageMem)), extractMemorySize(string(gpuMem)), nil
}

// IsVcgencmdInstalled verifica se o comando vcgencmd está instalado no sistema.
func IsVcgencmdInstalled() bool {
	_, err := exec.LookPath("vcgencmd")
	return err == nil
}

// GetDeviceName retorna o nome do dispositivo baseado na revisão do CPU.
func GetDeviceName() (string, error) {
	revision, err := GetCPURevision()
	if err != nil {
		return "", err
	}

	if deviceName, exists := deviceRevisions[revision]; exists {
		return deviceName, nil
	}
	return "", errors.New("device not recognized")
}
