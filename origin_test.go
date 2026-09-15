package main

import "testing"

func TestParseOriginBuildsThePrefixTheBrowserWillSend(t *testing.T) {
	for line, want := range map[string]string{
		"http://a.example.com/x":  "http://a.example.com/",
		"https://a.example.com/x": "https://a.example.com/",
		// The URL Standard drops a default port, so writing one must not
		// produce a prefix the browser can never match.
		"https://a.example.com:443/x": "https://a.example.com/",
		"http://a.example.com:80/x":   "http://a.example.com/",
		"https://c.example.com:8443/": "https://c.example.com:8443/",
		"HTTPS://D.Example.COM/x":     "https://d.example.com/",
		"http://192.0.2.9:65432/":     "http://192.0.2.9:65432/",
	} {
		got, err := parseOrigin(line)
		if err != nil {
			t.Errorf("parseOrigin(%q): %v", line, err)
			continue
		}
		if got.Prefix != want {
			t.Errorf("parseOrigin(%q).Prefix = %q, want %q", line, got.Prefix, want)
		}
	}
}

func TestParseOriginRejectsLinesThatCouldNeverMatch(t *testing.T) {
	for _, line := range []string{
		"",
		"   ",
		"http://a.example.com/x ",
		"a.example.com",
		"a.example.com:8443",
		"192.0.2.9",
		"ftp://a.example.com/x",
		"http://",
		"https:///x",
		"http://example.com./x",
		"https://[2001:db8::1]/",
		"http://2001:db8::1/",
		"http://*.example.com/",
		"http://.example.com/",
		"http://a..b.example.com/",
	} {
		if got, err := parseOrigin(line); err == nil {
			t.Errorf("parseOrigin(%q) = %+v, want an error", line, got)
		}
	}
}
