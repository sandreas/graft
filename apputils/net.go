package apputils

import (
	"net"
)

func GetOutboundIpAsString(fallbackValue string , dialCallback func (network, address string) (net.Conn, error)) (string, error) {
	conn, err := dialCallback("udp", "8.8.8.8:80")
	if err != nil {
		return fallbackValue, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
