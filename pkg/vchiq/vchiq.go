// Package vchiq provides functions to retrieve various system information
// using vcgencmd and other system interfaces on Raspberry Pi.
package vchiq

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Constants for vcgencmd get_throttled result
const (
	UnderVoltage         = 1 << iota // 1
	FreqCap                          // 2
	Throttling                       // 4
	UnderVoltageOccurred             // 65536
	FreqCapOccurred                  // 131072
	Throttled                        // 262144
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
	return string(usageMem), string(gpuMem), nil
}

// GetCPUTemp returns the CPU temperature as a string.
func GetCPUTemp() (string, error) {
	temp, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "", errors.New("Permission Denied")
	}
	cpuTemp, err := strconv.ParseInt(strings.TrimSpace(string(temp[:5])), 10, 64)
	if err != nil {
		return "", errors.New("Error converting to int")
	}
	cpuTempC := float64(cpuTemp) / 1000.0
	return strconv.FormatFloat(cpuTempC, 'f', 2, 64) + "C", nil
}

// GetLoadAverage returns the 1-minute load average as a string.
func GetLoadAverage() (string, error) {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		return "", err
	}
	uptimeResult := string(out)
	loadIndex := strings.Index(uptimeResult, "load average:")
	loadValue := uptimeResult[loadIndex+len("load average:"):]
	load := strings.TrimSpace(strings.Split(loadValue, ",")[0])
	return load, nil
}

// GetMemoryUsagePercent returns the system memory usage as a percentage.
func GetMemoryUsagePercent() (float64, error) {
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return 0.0, err
	}
	used := sysInfo.Totalram - sysInfo.Freeram
	usedPercent := float64(used) / float64(sysInfo.Totalram)
	return usedPercent, nil
}

// GetMemory returns the total, free, and used memory in bytes.
func GetMemory() (float64, float64, float64, error) {
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return 0, 0, 0, err
	}
	total := float64(sysInfo.Totalram)
	free := float64(sysInfo.Freeram)
	used := total - free
	return total, free, used, nil
}

// GetDiskUsage returns the usage percentage of boot, root, and home partitions.
func GetDiskUsage() (float64, float64, float64, error) {
	bootUsage, err := getDiskUsage("/boot")
	if err != nil {
		return 0, 0, 0, err
	}
	rootUsage, err := getDiskUsage("/")
	if err != nil {
		return 0, 0, 0, err
	}
	homeUsage, err := getDiskUsage("/home")
	if err != nil {
		return 0, 0, 0, err
	}
	return bootUsage, rootUsage, homeUsage, nil
}

// getDiskUsage returns the usage percentage of the specified path.
func getDiskUsage(path string) (float64, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, err
	}
	total := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := total - free
	usedPercent := float64(used) / float64(total)
	return usedPercent, nil
}

// clean removes specified substrings from the input string and trims the result.
func clean(str string, args ...string) string {
	for _, arg := range args {
		str = strings.ReplaceAll(str, arg, "")
	}
	return strings.TrimSpace(str)
}

// IsVcgencmdInstalled verifica se o comando vcgencmd está instalado no sistema.
func IsVcgencmdInstalled() bool {
	_, err := exec.LookPath("vcgencmd")
	return err == nil
}

func GetHostname() (string, error) {
	// Tenta obter o hostname usando o comando 'hostname'
	host, err := exec.Command("hostname").Output()
	if err == nil {
		return clean(string(host)), nil
	}

	// Se falhar, tenta ler o hostname do arquivo '/etc/hostname'
	hostFile, err := exec.Command("cat", "/etc/hostname").Output()
	if err != nil {
		return "", errors.New("couldn't get hostname from both 'hostname' command and /etc/hostname file")
	}

	return clean(string(hostFile)), nil
}
