package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/jessevdk/go-flags"
)

type option struct {
	Input     flags.Filename    `short:"i" long:"input" description:"input template file"`
	Output    flags.Filename    `short:"o" long:"output" description:"output file"`
	Values    map[string]string `short:"v" long:"values" description:"specifies key:value pairs to be loaded into the template"`
	Keys      []string          `short:"k" long:"keys" description:"specifies environmental variables to be loaded into the template"`
	Overwrite bool              `short:"f" long:"force-overwrite" description:"overwrite output file if exists"`
	Quiet     bool              `long:"quiet" description:"Quiet mode disables non-error outputs"`
    Version   bool              `long:"version" description:"Output version tag and exit"`
}

func (o option) Usage() string {
    return "fcfg --input=openvpn.conf.tmpl --output=openvpn.conf --values=ca:ca.crt --values=key:client.key"
}

var version string

func main() {
	o := option{
		Values: make(map[string]string),
	}

	if _, err := flags.Parse(&o); err != nil {
		os.Exit(-1)
	}

    if o.Version {
        fmt.Printf("%s\n", version)
        os.Exit(0)
    }

	if !o.Quiet {
		fmt.Printf("ryankurte/utils fcfg version: %s\n", version)
		fmt.Printf("https://github.com/ryankurte/utils\n")
    }

    if o.Input == "" || o.Output == "" {
        fmt.Printf("Missing input template file (-i, --input) and/or output file (-o, --output) arguments\n")
        os.Exit(-2)
    }

    if !o.Quiet {
        fmt.Printf("Loading template file: %s\n", o.Input)
	}

	f, err := ioutil.ReadFile(string(o.Input))
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		os.Exit(-2)
	}

	tmpl, err := template.New("").Parse(string(f))
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		os.Exit(-2)
	}

	values := make(map[string]string)
	for _, k := range o.Keys {
		values[k] = os.Getenv(k)
	}
	for k, v := range o.Values {
		values[k] = v
	}

	if !o.Quiet {
		fmt.Printf("Loaded values:\n")
		for k, v := range values {
			fmt.Printf("  - %s:%s\n", k, v)
		}
	}

	if !o.Quiet {
		fmt.Printf("Writing output file: %s\n", o.Output)
	}

	wr, err := os.Create(string(o.Output))
	if err != nil {
		fmt.Printf("Error creating output file: %s\n", err)
		os.Exit(-2)
	}

	w := bufio.NewWriter(wr)
	err = tmpl.Execute(w, values)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		os.Exit(-3)
	}

	w.Flush()
	wr.Close()

    if !o.Quiet {
        fmt.Printf("Output file configured\n")
    }
}
