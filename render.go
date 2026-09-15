package main

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed pac.tmpl
var pacTemplate string

// Not the raw CIDRs: deciding what to emit inside the template would put that
// out of reach of the tests.
type pacData struct {
	Proxy    string
	Networks []network
	Covered  []string
	Origins  []hostOrigins
}

var tmpl = template.Must(template.New("pac").Parse(pacTemplate))

func render(w io.Writer, d pacData) error {
	return tmpl.Execute(w, d)
}
