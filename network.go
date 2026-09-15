package main

import (
	"fmt"
	"net"
)

type network struct {
	CIDR    string
	Address string
	Mask    string
}

func parseNet(s string) (*net.IPNet, error) {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("not a valid CIDR: %q", s)
	}
	// Not accepted: isInNet() handles IPv6 differently across browsers.
	if len(n.Mask) != net.IPv4len {
		return nil, fmt.Errorf("IPv6 is not supported: %q", s)
	}
	return n, nil
}

func toNetworks(nets []*net.IPNet) []network {
	out := make([]network, 0, len(nets))
	for _, n := range nets {
		out = append(out, network{
			CIDR:    n.String(),
			Address: n.IP.String(),
			Mask:    net.IP(n.Mask).String(),
		})
	}
	return out
}
