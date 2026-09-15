package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	err := run()
	if err == nil || errors.Is(err, errDone) {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// Not spread across the package: I/O here keeps the rest pure.
func run() error {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		return err
	}
	origins, err := readOriginFile(opts.input)
	if err != nil {
		return err
	}
	pac, err := buildPAC(origins, opts.nets, opts.proxy)
	if err != nil {
		return err
	}
	return writePAC(opts.output, pac)
}

// Not stdin: a pipe hides the exit status of whatever produced the list.
func readOriginFile(path string) ([]origin, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readOrigins(f)
}

// Not streamed: a template failure would truncate a working PAC.
func writePAC(path string, b []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(b)
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
