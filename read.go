package main

import (
	"bufio"
	"fmt"
	"io"
)

// Not skipped: a line-oriented extraction writes an empty value as an empty
// line, so skipping would hide an entry that dropped out upstream.
func readOrigins(r io.Reader) ([]origin, error) {
	var origins []origin
	sc := bufio.NewScanner(r)
	for line := 1; sc.Scan(); line++ {
		o, err := parseOrigin(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		origins = append(origins, o)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return origins, nil
}
