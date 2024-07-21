package metricagent

import (
	"errors"
	"fmt"
	"net"
)

func identifyIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to identify IP addresses: %w", err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.To4(), nil
		}
	}
	return nil, errors.New("ip address is not found")
}
