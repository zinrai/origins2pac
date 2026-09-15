package main

import "testing"

func TestToNetworksGivesIsInNetTheAddressAndMaskItExpects(t *testing.T) {
	for cidr, want := range map[string]network{
		"192.0.2.0/24":      {CIDR: "192.0.2.0/24", Address: "192.0.2.0", Mask: "255.255.255.0"},
		"198.51.100.128/25": {CIDR: "198.51.100.128/25", Address: "198.51.100.128", Mask: "255.255.255.128"},
		"203.0.113.0/26":    {CIDR: "203.0.113.0/26", Address: "203.0.113.0", Mask: "255.255.255.192"},
		"198.51.100.7/32":   {CIDR: "198.51.100.7/32", Address: "198.51.100.7", Mask: "255.255.255.255"},
		"0.0.0.0/0":         {CIDR: "0.0.0.0/0", Address: "0.0.0.0", Mask: "0.0.0.0"},
		// Not the address as typed: isInNet() would judge against a host.
		"203.0.113.42/24": {CIDR: "203.0.113.0/24", Address: "203.0.113.0", Mask: "255.255.255.0"},
	} {
		got := toNetworks(mustParseNets(t, cidr))
		if len(got) != 1 || got[0] != want {
			t.Errorf("toNetworks(%q) = %+v, want %+v", cidr, got, want)
		}
	}
}

func TestParseNetRejectsWhatIsInNetCannotJudge(t *testing.T) {
	for _, cidr := range []string{
		"2001:db8::/32",
		"192.0.2.0",
		"192.0.2.0/33",
		"",
	} {
		if _, err := parseNet(cidr); err == nil {
			t.Errorf("parseNet(%q) = nil error, want an error", cidr)
		}
	}
}
