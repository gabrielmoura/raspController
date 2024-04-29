package mdns // import "example.com/mdns"

import (
	"github.com/hashicorp/mdns"
	"net"
	"os"
)

// getLocalIP retorna o endereço IP local IPv4 não loopback, se disponível.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

// SetDNS configura o serviço de DNS multicast (mDNS).
func SetDNS(appName string, appPort int) error {
	// Obtém o nome do host
	host, err := os.Hostname()
	if err != nil {
		return err
	}

	// Cria as informações do serviço
	info := []string{appName}

	// Cria e configura o serviço mDNS
	service, err := mdns.NewMDNSService(host, "_rpi._tcp", "", getLocalIP(), appPort, nil, info)
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
