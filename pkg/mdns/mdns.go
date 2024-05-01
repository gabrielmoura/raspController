package mdns

import (
	"net"
	"os"

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
	host, _ := os.Hostname()

	// Cria as informações do serviço
	info := []string{appName}

	// Cria e configura o serviço mDNS
	service, err := mdns.NewMDNSService(host, "_rpi._tcp", "", "", appPort, getLocalIP(), info)
	if err != nil {
		return err
	}

	// Cria o servidor mDNS e defer sua parada
	_, err = mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}

	return nil
}
