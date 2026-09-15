# origins2pac

Generate a PAC (Proxy Auto-Configuration) file from a list of origins and a set of IPv4 CIDRs.

origins2pac generates it from a plain list, one URL per line. The result is a static file.

## Usage

```
origins2pac -input FILE [flags]
```

| Flag | Meaning |
|---|---|
| `-input FILE` | Origin list to read. Required |
| `-net CIDR` | IPv4 CIDR whose addresses go through the proxy. May be repeated |
| `-output FILE` | Write the PAC here instead of standard output |
| `-proxy STRING` | Proxy string the PAC returns. Defaults to a local proxy on port 3128 |
| `-version` | Print the version and exit |

The proxy string is written out unchanged, so `PROXY proxy.example.com:3128; DIRECT` or a `SOCKS5` entry works too.

origins2pac exits 0 on success and 1 on any failure. Errors go to standard error as a single line.

## Input

One URL per line, and nothing else. No comments, no blank lines, no surrounding whitespace. origins2pac reads a file, not standard input.

Every line starts with `http://` or `https://`. What decides the routing is the origin, which is the scheme, the host and the port. The path is kept as a comment in the generated file and takes no part in the matching.

```
http://a.example.com/app-a/
http://a.example.com/app-b/
https://b.example.com/x
https://c.example.com:8443/
http://203.0.113.5/
http://192.0.2.9:65432/
```

A line without a port covers the default port for its scheme, 80 for `http` and 443 for `https`. Writing that port out is the same thing. A line with any other port covers that port alone, so `https://c.example.com:8443/` leaves `https://c.example.com/` direct. To cover a host on both schemes, list it twice.

Host names are lowercased. Wildcards, suffix forms, IPv6 literals and trailing dots are rejected.

Anything else fails, naming the line:

```
Error: line 4: every line must start with http:// or https://: "d.example.com"
```

## Output

```
$ origins2pac -input origins.txt -net 192.0.2.0/24 -proxy 'PROXY proxy.example.com:3128'
function FindProxyForURL(url, host) {
    var proxy = "PROXY proxy.example.com:3128";

    // Covered by the networks below:
    //   http://192.0.2.9:65432/

    // 192.0.2.0/24
    if (isInNet(host, "192.0.2.0", "255.255.255.0")) {
        return proxy;
    }

    var proxyOrigins = {
        "203.0.113.5": [
            "http://203.0.113.5/"
        ],
        // /app-a/
        // /app-b/
        "a.example.com": [
            "http://a.example.com/"
        ],
        // /x
        "b.example.com": [
            "https://b.example.com/"
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
```

The same list always produces the same bytes.

An empty list is an error, whether or not `-net` is given:

```
Error: the origin list is empty
```

## License

This project is licensed under the [MIT License](LICENSE).
