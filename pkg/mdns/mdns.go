package mdns

import (
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

// getLocalIP returns the non-loopback IPv4 local IP address, if available.
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

// SetDNS configures the multicast DNS (mDNS) service.
func SetDNS(appName string, appPort int) error {
	// Get the host name
	host, _ := os.Hostname()

	// Creates the service information
	info := []string{appName}

	// Creates and configures the mDNS service
	service, err := mdns.NewMDNSService(host, "_rpi._tcp", "", "", appPort, getLocalIP(), info)
	if err != nil {
		return err
	}

	// Create mDNS server and defer its stop
	_, err = mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}

	return nil
}
