package domain

import (
	"fmt"
	"net"
	"strconv"
)

type Address struct {
	host string
	port int
}

func NewAddressFromHost(host string) (*Address, error) {
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		return nil, InvalidParameterError("host")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, InvalidParameterError("port")
	}

	return NewAddress(host, portInt)
}

func NewAddressesFromHosts(hosts []string) ([]*Address, error) {
	var ret []*Address

	for _, host := range hosts {
		address, err := NewAddressFromHost(host)
		if err != nil {
			return nil, err
		}

		ret = append(ret, address)
	}

	return ret, nil
}

func NewAddress(host string, port int) (*Address, error) {
	if host == "" {
		return nil, InvalidParameterError("host")
	}

	if net.ParseIP(host) == nil {
		return nil, InvalidParameterError("host")
	}

	if port < 1 || port > 65535 {
		return nil, InvalidParameterError("port")
	}

	return &Address{
		host: host,
		port: port,
	}, nil
}

func (a *Address) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}
