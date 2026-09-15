package main

import (
	"bytes"
	"testing"
)

const goldenPAC = `function FindProxyForURL(url, host) {
    var proxy = "PROXY proxy.example.com:3128";

    // Covered by the networks below:
    //   http://198.51.100.200/x

    // 192.0.2.0/24
    if (isInNet(host, "192.0.2.0", "255.255.255.0")) {
        return proxy;
    }

    // 198.51.100.128/25
    if (isInNet(host, "198.51.100.128", "255.255.255.128")) {
        return proxy;
    }

    var proxyOrigins = {
        // /app-a/
        // /app-b/
        "a.example.com": [
            "http://a.example.com/",
            "https://a.example.com/"
        ],
        "c.example.com": [
            "https://c.example.com:8443/"
        ]
    };

    var prefixes = proxyOrigins[host];
    if (prefixes) {
        for (var i = 0; i < prefixes.length; i++) {
            if (url.indexOf(prefixes[i]) === 0) {
                return proxy;
            }
        }
    }

    return "DIRECT";
}
`

// Not asserted field by field: a broken template still renders, and the
// browser then sends every origin to DIRECT in silence.
func TestRenderMatchesTheGoldenPAC(t *testing.T) {
	nets := mustParseNets(t, "192.0.2.0/24", "198.51.100.128/25")

	var got bytes.Buffer
	err := render(&got, pacData{
		Proxy:    "PROXY proxy.example.com:3128",
		Networks: toNetworks(nets),
		Covered:  []string{"http://198.51.100.200/x"},
		Origins: []hostOrigins{
			{
				Host:  "a.example.com",
				Paths: []string{"/app-a/", "/app-b/"},
				Prefixes: []string{
					"http://a.example.com/",
					"https://a.example.com/",
				},
			},
			{
				Host:     "c.example.com",
				Prefixes: []string{"https://c.example.com:8443/"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got.String() != goldenPAC {
		t.Errorf("rendered PAC differs\n--- got ---\n%s\n--- want ---\n%s", got.String(), goldenPAC)
	}
}
