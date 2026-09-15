package main

import (
	"strings"
	"testing"
)

func TestBuildPACJudgesOriginsWithTheNetsItEmits(t *testing.T) {
	pac, err := buildPAC(
		[]origin{
			{Host: "192.0.2.9", Prefix: "http://192.0.2.9:65432/"},
			{Host: "a.example.com", Prefix: "https://a.example.com/"},
		},
		mustParseNets(t, "192.0.2.0/24"),
		"PROXY proxy.example.com:3128")
	if err != nil {
		t.Fatal(err)
	}

	got := string(pac)
	if !strings.Contains(got, `isInNet(host, "192.0.2.0", "255.255.255.0")`) {
		t.Errorf("no isInNet() for the CIDR it was given:\n%s", got)
	}
	if strings.Contains(got, "192.0.2.9") {
		t.Errorf("lists an address its own isInNet() already covers:\n%s", got)
	}
	if !strings.Contains(got, `"https://a.example.com/"`) {
		t.Errorf("dropped an origin no CIDR covers:\n%s", got)
	}
}

func TestBuildPACRefusesAnEmptyOriginList(t *testing.T) {
	if _, err := buildPAC(nil, nil, "PROXY proxy.example.com:3128"); err == nil {
		t.Error("an empty list = nil error, want an error")
	}
	if _, err := buildPAC(nil, mustParseNets(t, "192.0.2.0/24"), "PROXY proxy.example.com:3128"); err == nil {
		t.Error("an empty list alongside a CIDR = nil error, want an error")
	}
}

func TestBuildPACKeepsAListTheNetsWhollyCover(t *testing.T) {
	pac, err := buildPAC(
		[]origin{{
			Host:   "192.0.2.9",
			Path:   "/x",
			Line:   "http://192.0.2.9/x",
			Prefix: "http://192.0.2.9/",
		}},
		mustParseNets(t, "192.0.2.0/24"),
		"PROXY proxy.example.com:3128")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pac), "http://192.0.2.9/x") {
		t.Errorf("the covered origin is named nowhere:\n%s", pac)
	}
}
