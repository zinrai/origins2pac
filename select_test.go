package main

import (
	"net"
	"slices"
	"testing"
)

func TestSelectOriginsMergesByHostAndDropsWhatIsInNetCovers(t *testing.T) {
	got, covered := selectOrigins([]origin{
		{Host: "z.example.com", Path: "/", Line: "https://z.example.com/",
			Prefix: "https://z.example.com/"},
		{Host: "192.0.2.9", Path: "/", Line: "http://192.0.2.9:65432/",
			Prefix: "http://192.0.2.9:65432/"},
		{Host: "a.example.com", Path: "/app-b/", Line: "https://a.example.com/app-b/",
			Prefix: "https://a.example.com/"},
		{Host: "a.example.com", Path: "/app-a/", Line: "http://a.example.com/app-a/",
			Prefix: "http://a.example.com/"},
		{Host: "a.example.com", Path: "/app-b/", Line: "https://a.example.com/app-b/",
			Prefix: "https://a.example.com/"},
	}, mustParseNets(t, "192.0.2.0/24"))

	want := []hostOrigins{
		{
			Host:     "a.example.com",
			Paths:    []string{"/app-a/", "/app-b/"},
			Prefixes: []string{"http://a.example.com/", "https://a.example.com/"},
		},
		{
			// A path of "/" says nothing the host key does not already say.
			Host:     "z.example.com",
			Prefixes: []string{"https://z.example.com/"},
		},
	}
	if !slices.EqualFunc(got, want, func(a, b hostOrigins) bool {
		return a.Host == b.Host &&
			slices.Equal(a.Paths, b.Paths) &&
			slices.Equal(a.Prefixes, b.Prefixes)
	}) {
		t.Errorf("selectOrigins() = %+v, want %+v", got, want)
	}

	// Not silently gone: the file would otherwise lose the address entirely.
	if !slices.Equal(covered, []string{"http://192.0.2.9:65432/"}) {
		t.Errorf("selectOrigins() covered = %q, want the address the CIDR swallowed", covered)
	}
}

func mustParseNets(t *testing.T, cidrs ...string) []*net.IPNet {
	t.Helper()
	var nets []*net.IPNet
	for _, c := range cidrs {
		n, err := parseNet(c)
		if err != nil {
			t.Fatalf("parseNet(%q): %v", c, err)
		}
		nets = append(nets, n)
	}
	return nets
}
