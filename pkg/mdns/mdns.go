package mdns

import (
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/mdns"
)

// getLocalIP retorna o endereço IP local IPv4 não loopback, se disponível.
func getLocalIP() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
			continue
		}
		return []net.IP{ipNet.IP}
	}
	return nil

}

// SetDNS configura o serviço de DNS multicast (mDNS).
func SetDNS(appName string, appPort int) error {
	// Obtém o nome do host
	host, err := getHostname()
	if err != nil {
		return err
	}

	// Cria as informações do serviço
	info := []string{appName}

	// Cria e configura o serviço mDNS
	service, err := mdns.NewMDNSService(host, "_rpi._tcp", "", host, appPort, getLocalIP(), info)
	if err != nil {
		return err
	}

	// Cria o servidor mDNS e defer sua parada
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	defer func() {
		_ = server.Shutdown()
	}()

	return nil
}

func getHostname() (string, error) {
	host, err := os.Hostname()
	if err == nil {
		// Adiciona um ponto ao final do nome do host para torná-lo um FQDN válido
		return clean(host) + ".", nil
	}

	hostFile, err := exec.Command("cat", "/etc/hostname").Output()
	if err != nil {
		return "", err
	}

	// Adiciona um ponto ao final do nome do host lido do arquivo '/etc/hostname' para torná-lo um FQDN válido
	return clean(string(hostFile)) + ".", nil
}

// clean removes specified substrings from the input string and trims the result.
func clean(str string, args ...string) string {
	for _, arg := range args {
		str = strings.ReplaceAll(str, arg, "")
	}
	return strings.TrimSpace(str)
}
