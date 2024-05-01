// Package vchiq provides functions to retrieve various system information
// using vcgencmd and other system interfaces on Raspberry Pi.
package vchiq

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// Constants for vcgencmd get_throttled result
const (
	UnderVoltage        int64 = 1
	FreqCap                   = 1 << 1
	Throttling                = 1 << 2
	UnderVoltageOccured       = 1 << 16
	FreqCapOccured            = 1 << 17
	Throttled                 = 1 << 18
)

// GetThrottled returns the throttled status as an integer.
func GetThrottled() (int64, error) {
	rawThrottled, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return 0, errors.New("couldn't run vcgencmd")
	}
	throttled, err := strconv.ParseInt(strings.TrimSpace(string(rawThrottled[12:])), 16, 32)
	if err != nil {
		return 0, errors.New("couldn't parse throttled output")
	}
	return throttled, nil
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

// GetDeviceName retorna o nome do dispositivo com base no código de revisão da CPU.
// Consulte https://www.raspberrypi.org/documentation/hardware/raspberrypi/revision-codes/README.md
func GetDeviceName() (string, error) {
	revision, err := GetCPURevision()
	if err != nil {
		return "", err
	}

	switch revision {
	case "900021":
		return "Raspberry Pi Zero", nil
	case "900032":
		return "Raspberry Pi Zero W", nil
	case "900092", "920092":
		return "Raspberry Pi Zero", nil
	case "900093", "920093":
		return "Raspberry Pi Zero", nil
	case "9000c1":
		return "Raspberry Pi Zero W", nil
	case "9020e0", "9020e1":
		return "Raspberry Pi 3 Model A+", nil
	case "a01040", "a01041", "a21041":
		return "Raspberry Pi 2 Model B", nil
	case "a02082", "a22082", "a32082", "a52082":
		return "Raspberry Pi 3 Model B", nil
	case "a020a0", "a220a0":
		return "Raspberry Pi Compute Module 3", nil
	case "a020d3", "a22083":
		return "Raspberry Pi 3 Model B+", nil
	case "a020d4":
		return "Raspberry Pi 3 Model B+", nil
	case "a02042", "a22042":
		return "Raspberry Pi 2 Model B (with BCM2837)", nil
	case "a02100":
		return "Raspberry Pi Compute Module 3+", nil
	case "a03111", "b03111", "c03111":
		return "Raspberry Pi 4 Model B", nil
	case "b03112", "c03112":
		return "Raspberry Pi 4 Model B", nil
	case "b03114", "c03114":
		return "Raspberry Pi 4 Model B", nil
	case "b03115", "c03115":
		return "Raspberry Pi 4 Model B", nil
	case "c03130":
		return "Raspberry Pi 400", nil
	case "a03140", "b03140", "c03140", "d03140":
		return "Raspberry Pi Compute Module 4", nil
	case "902120":
		return "Raspberry Pi Zero 2 W", nil
	case "c04170":
		return "Raspberry Pi 5", nil
	case "d04170":
		return "Raspberry Pi 5", nil
	default:
		return "", errors.New("device not recognized")
	}
}
