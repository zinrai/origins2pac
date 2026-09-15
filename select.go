package main

import (
	"cmp"
	"net"
	"slices"
)

type hostOrigins struct {
	Host     string
	Paths    []string
	Prefixes []string
}

// Not left to the PAC: a covered address duplicates what isInNet() already
// decides. Not dropped either: it would otherwise vanish from the file.
func selectOrigins(origins []origin, nets []*net.IPNet) ([]hostOrigins, []string) {
	var covered []string
	prefixes := make(map[string][]string)
	paths := make(map[string][]string)
	for _, o := range origins {
		if coveredByNets(o.Host, nets) {
			covered = append(covered, o.Line)
			continue
		}
		prefixes[o.Host] = append(prefixes[o.Host], o.Prefix)
		if o.Path != "/" {
			paths[o.Host] = append(paths[o.Host], o.Path)
		}
	}

	out := make([]hostOrigins, 0, len(prefixes))
	for host, p := range prefixes {
		slices.Sort(p)
		q := paths[host]
		slices.Sort(q)
		out = append(out, hostOrigins{
			Host:     host,
			Paths:    slices.Compact(q),
			Prefixes: slices.Compact(p),
		})
	}
	slices.SortFunc(out, func(a, b hostOrigins) int { return cmp.Compare(a.Host, b.Host) })

	slices.Sort(covered)
	return out, slices.Compact(covered)
}

func coveredByNets(host string, nets []*net.IPNet) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
