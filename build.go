package main

import (
	"bytes"
	"fmt"
	"net"
)

// Not inlined into run(): the CIDRs must reach both the isInNet() blocks and
// the origin selection, which has to be checkable without touching files.
func buildPAC(origins []origin, nets []*net.IPNet, proxy string) ([]byte, error) {
	// Not checked after selection: a list the CIDRs wholly cover still said
	// something, and -net does not stand in for a list that said nothing.
	if len(origins) == 0 {
		return nil, fmt.Errorf("the origin list is empty")
	}
	selected, covered := selectOrigins(origins, nets)

	var buf bytes.Buffer
	err := render(&buf, pacData{
		Proxy:    proxy,
		Networks: toNetworks(nets),
		Covered:  covered,
		Origins:  selected,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
