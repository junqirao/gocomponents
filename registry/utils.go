package registry

import (
	"net"
)

func getIp() (v4 net.IP, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	var addresses []net.Addr
	for _, inf := range interfaces {
		if inf.Flags|net.FlagUp > 0 && inf.Flags|net.FlagRunning > 0 {
			addresses, err = inf.Addrs()
			if err != nil {
				return
			}
			for _, addr := range addresses {
				if ipAddr, ok := addr.(*net.IPNet); ok &&
					ipAddr.IP != nil &&
					ipAddr.IP.To4() != nil &&
					!ipAddr.IP.IsLoopback() {
					v4 = ipAddr.IP
					return
				}
			}
		}
	}

	return
}
