package vchiq

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// GetCPUTemp returns the CPU temperature as a string.
func GetCPUTemp() (string, error) {
	temp, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "", errors.New("permission Denied")
	}
	cpuTemp, err := strconv.ParseInt(strings.TrimSpace(string(temp[:5])), 10, 64)
	if err != nil {
		return "", errors.New("error converting to int")
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

// RetrieveDiskUsagePercent returns the usage percentage of boot, root, and home partitions.
func RetrieveDiskUsagePercent() (float64, float64, float64, error) {
	bootUsage, err := calculateDiskUsage("/boot", true)
	if err != nil {
		return 0, 0, 0, err
	}
	rootUsage, err := calculateDiskUsage("/", true)
	if err != nil {
		return 0, 0, 0, err
	}
	homeUsage, err := calculateDiskUsage("/home", true)
	if err != nil {
		return 0, 0, 0, err
	}
	return bootUsage, rootUsage, homeUsage, nil
}

// RetrieveDiskUsage returns the usage of boot, root, and home partitions in bytes.
func RetrieveDiskUsage() (int, int, int, error) {
	bootUsage, err := calculateDiskUsage("/boot", false)
	if err != nil {
		return 0, 0, 0, err
	}
	rootUsage, err := calculateDiskUsage("/", false)
	if err != nil {
		return 0, 0, 0, err
	}
	homeUsage, err := calculateDiskUsage("/home", false)
	if err != nil {
		return 0, 0, 0, err
	}
	return int(bootUsage), int(rootUsage), int(homeUsage), nil
}

// RetrieveDiskTotal returns the total disk space of boot, root, and home partitions in bytes.
func RetrieveDiskTotal() (int, int, int, error) {
	bootTotal, _, _, err := getDiskSize("/boot")
	if err != nil {
		return 0, 0, 0, err
	}
	rootTotal, _, _, err := getDiskSize("/")
	if err != nil {
		return 0, 0, 0, err
	}
	homeTotal, _, _, err := getDiskSize("/home")
	if err != nil {
		return 0, 0, 0, err
	}
	return int(bootTotal), int(rootTotal), int(homeTotal), nil
}

func calculateDiskUsage(path string, percent bool) (float64, error) {
	var fs syscall.Statfs_t

	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, err
	}

	totalBlocks := fs.Blocks
	blockSize := fs.Bsize
	totalSpace := float64(totalBlocks * uint64(blockSize))
	freeSpace := float64(fs.Bfree * uint64(blockSize))
	usedSpace := totalSpace - freeSpace
	usedPercent := usedSpace / totalSpace

	if percent {
		return usedPercent, nil
	}

	return usedSpace, nil
}

func getDiskSize(path string) (float64, float64, float64, error) {
	var fs syscall.Statfs_t

	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, 0, 0, err
	}
	totalSpace := float64(fs.Blocks * uint64(fs.Bsize))
	freeSpace := float64(fs.Bfree * uint64(fs.Bsize))
	availableSpace := float64(fs.Bavail * uint64(fs.Bsize))
	return totalSpace, freeSpace, availableSpace, nil
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

// getCPUInfoValue busca um valor específico no arquivo cpuinfo com base no prefixo fornecido.
func getCPUInfoValue(prefix string) (string, error) {
	// Lê o conteúdo do arquivo cpuinfo
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	// Converte os dados do cpuinfo para uma string
	cpuInfoStr := string(cpuInfo)

	// Procura o campo correspondente ao prefixo
	lines := strings.Split(cpuInfoStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				return value, nil
			}
		}
	}

	return "", errors.New("couldn't find value for prefix in /proc/cpuinfo")
}

// GetCPUSerial retorna o código de serial da CPU.
func GetCPUSerial() (string, error) {
	return getCPUInfoValue("Serial")
}

// GetCPURevision retorna o código de revisão da CPU.
func GetCPURevision() (string, error) {
	return getCPUInfoValue("Revision")
}

func GetCPUModel() (string, error) {
	return getCPUInfoValue("model name")
}
func GetCPUCores() (string, error) {
	return getCPUInfoValue("cpu cores")
}

func GetCPUMhz() (string, error) {
	return getCPUInfoValue("cpu MHz")
}

// GetUptime retorna o tempo de atividade do sistema.
func GetUptime() (string, error) {
	uptime, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}

	uptimeStr := strings.Split(string(uptime), " ")[0]
	uptimeInt, err := strconv.ParseFloat(uptimeStr, 64)
	if err != nil {
		return "", err
	}

	uptimeDuration := time.Duration(int64(uptimeInt)) * time.Second
	return uptimeDuration.String(), nil
}

// GetKernelVersion retorna a versão do kernel.
func GetKernelVersion() (string, error) {
	// Lê o conteúdo do arquivo /proc/version
	version, err := os.ReadFile("/proc/version")
	if err != nil {
		return "", err
	}

	// Converte os dados do arquivo para uma string
	versionStr := string(version)

	// Divide a string por espaços
	parts := strings.Fields(versionStr)

	// Retorna a segunda parte da string
	return parts[2], nil
}

// GetFqdn retorna o nome de domínio totalmente qualificado do sistema.
func GetFqdn() (string, error) {
	// Tenta obter o FQDN usando o comando 'hostname'
	fqdn, err := exec.Command("hostname", "-f").Output()
	if err == nil {
		return clean(string(fqdn)), nil
	}

	// Se falhar, tenta ler o FQDN do arquivo '/etc/hostname'
	fqdnFile, err := exec.Command("cat", "/etc/hostname").Output()
	if err != nil {
		return "", errors.New("couldn't get FQDN from both 'hostname -f' command and /etc/hostname file")
	}

	return clean(string(fqdnFile)), nil
}

// GetLocalIP retorna um slice de endereços IPs.
func GetIps() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
			continue
		}
		return []net.IP{ipNet.IP}, nil
	}
	return nil, errors.New("couldn't get local IP")

}

// GetOsName retorna o nome do sistema operacional.
func GetOsName() (string, error) {
	// pegue de /etc/os-release
	osRelease, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}

	// Converte os dados do arquivo para uma string

	osReleaseStr := string(osRelease)

	// Procura o campo de nome do sistema operacional
	lines := strings.Split(osReleaseStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				osName := strings.Trim(parts[1], "\"")
				return osName, nil
			}
		}
	}

	return "", errors.New("couldn't find OS name in /etc/os-release")
}
