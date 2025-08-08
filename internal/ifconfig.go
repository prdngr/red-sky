package internal

import (
	"errors"
	"io"
	"net"
	"net/http"
)

const ifconfigUrl = "https://ifconfig.me/ip"

func GetPublicIp() (net.IP, error) {
	response, err := http.Get(ifconfigUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(body))

	if ip == nil || ip.To4() == nil {
		return nil, errors.New("could not automatically determine IPv4 address")
	}

	return ip, nil
}
