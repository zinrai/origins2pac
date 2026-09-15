package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Not a field on options: -h and -version both end the run before anything is
// generated, and run() should not have to ask which one happened.
var errDone = errors.New("nothing left to do")

type options struct {
	input  string
	output string
	proxy  string
	nets   []*net.IPNet
}

// Not flag.ExitOnError: that exits 2 with a usage dump instead of one line.
func parseOptions(argv []string) (options, error) {
	fs := flag.NewFlagSet("origins2pac", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var opts options
	// Not a local copied back at the end: a lost assignment would drop every
	// -net without a word.
	fs.Var((*netFlag)(&opts.nets), "net", "IPv4 CIDR to match with isInNet(); repeatable")
	// Not a positional argument: flag parsing stops at the first one.
	fs.StringVar(&opts.input, "input", "", "origin list to read, one http or https URL per line")
	fs.StringVar(&opts.output, "output", "", "write the PAC to this file instead of stdout")
	fs.StringVar(&opts.proxy, "proxy", "PROXY 127.0.0.1:3128", "proxy string the PAC returns")
	showVersion := fs.Bool("version", false, "print the version and exit")

	if err := fs.Parse(argv); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(fs)
			return options{}, errDone
		}
		return options{}, err
	}
	if *showVersion {
		printVersion()
		return options{}, errDone
	}
	if fs.NArg() > 0 {
		return options{}, fmt.Errorf("unexpected argument %q; the origin list is given with -input", fs.Arg(0))
	}
	if opts.input == "" {
		return options{}, fmt.Errorf("-input is required")
	}
	return opts, nil
}

func printUsage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stdout, "Usage: origins2pac -input FILE [flags]")
	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()
}

type netFlag []*net.IPNet

func (f *netFlag) String() string {
	s := make([]string, 0, len(*f))
	for _, n := range *f {
		s = append(s, n.String())
	}
	return strings.Join(s, ",")
}

func (f *netFlag) Set(v string) error {
	n, err := parseNet(v)
	if err != nil {
		return err
	}
	*f = append(*f, n)
	return nil
}
