package utils

import (
	"errors"
	"net"
)

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", errors.New("net.UDPAddr type assertion failed")
	}
	return localAddr.IP.String(), nil
}
