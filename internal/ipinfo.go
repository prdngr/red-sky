package internal

import (
	"errors"
	"io"
	"net"
	"net/http"
)

const ipInfoUrl = "https://ipinfo.io/ip"

func GetPublicIp() (net.IP, *net.IPNet, error) {
	response, err := http.Get(ipInfoUrl)

	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, nil, err
	}

	ip, cidr, err := net.ParseCIDR(string(body) + "/32")

	if err != nil || ip == nil || ip.To4() == nil {
		return nil, nil, errors.New("could not automatically determine IPv4 address")
	}

	return ip, cidr, nil
}
