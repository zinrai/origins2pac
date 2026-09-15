package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"
)

type origin struct {
	Host   string
	Path   string
	Line   string
	Prefix string
}

func parseOrigin(line string) (origin, error) {
	if line == "" {
		return origin{}, fmt.Errorf("blank line")
	}
	if strings.ContainsFunc(line, unicode.IsSpace) {
		return origin{}, fmt.Errorf("contains whitespace: %q", line)
	}
	u, err := url.Parse(line)
	if err != nil {
		return origin{}, fmt.Errorf("not a valid URL: %q", line)
	}
	// Not any scheme: a PAC is consulted for http and https. Not a separate
	// check for a missing one either: net/url reads "a.example.com:8443" as a
	// scheme of its own, so both spellings land here.
	if u.Scheme != "http" && u.Scheme != "https" {
		return origin{}, fmt.Errorf("every line must start with http:// or https://: %q", line)
	}

	host := strings.ToLower(u.Hostname())
	if err := checkMatchable(host, line); err != nil {
		return origin{}, err
	}
	return origin{
		Host:   host,
		Path:   u.RequestURI(),
		Line:   line,
		Prefix: prefixFor(u.Scheme, host, u.Port()),
	}, nil
}

// Not kept when it is the default: the URL Standard sets the port to null when
// it matches the scheme, so ":80" never reaches the url a PAC is handed.
func prefixFor(scheme, host, port string) string {
	if port == "" || port == defaultPort(scheme) {
		return scheme + "://" + host + "/"
	}
	return scheme + "://" + host + ":" + port + "/"
}

func defaultPort(scheme string) string {
	if scheme == "https" {
		return "443"
	}
	return "80"
}

// Not repaired silently: a key that cannot equal the host argument sends its
// host to DIRECT with nothing reported.
func checkMatchable(host, line string) error {
	if host == "" {
		return fmt.Errorf("no host: %q", line)
	}
	if ip := net.ParseIP(host); ip != nil && ip.To4() == nil {
		return fmt.Errorf("IPv6 is not supported: %q", line)
	}
	if strings.HasSuffix(host, ".") {
		return fmt.Errorf("host must not end with a dot: %q", line)
	}
	if strings.Contains(host, "*") {
		return fmt.Errorf("wildcards are not supported: %q", line)
	}
	if strings.HasPrefix(host, ".") || strings.Contains(host, "..") {
		return fmt.Errorf("host has an empty label: %q", line)
	}
	return nil
}
