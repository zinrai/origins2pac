package main

import (
	"strings"
	"testing"
)

func TestReadOriginsKeepsEveryLine(t *testing.T) {
	got, err := readOrigins(strings.NewReader(
		"http://a.example.com/\nhttps://b.example.com/x\nhttps://c.example.com:8443/\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Errorf("readOrigins() returned %d origins, want 3: %+v", len(got), got)
	}
}

func TestReadOriginsStopsOnTheLineThatFails(t *testing.T) {
	_, err := readOrigins(strings.NewReader("http://a.example.com/\n\nhttp://b.example.com/\n"))
	if err == nil {
		t.Fatal("readOrigins() = nil error, want an error naming line 2")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("readOrigins() error = %v, want it to name line 2", err)
	}
}
