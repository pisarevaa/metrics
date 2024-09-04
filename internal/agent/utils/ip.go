package utils

import (
	"errors"
	"net"
)

const host = "8.8.8.8:80"

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", host)
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
