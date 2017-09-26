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
	Input     flags.Filename    `short:"i" long:"input" description:"input template file" required:"true"`
	Output    flags.Filename    `short:"o" long:"output" description:"output file" required:"true"`
	Values    map[string]string `short:"v" long:"values" description:"specifies key:value pairs to be loaded into the template"`
	Keys      []string          `short:"k" long:"keys" description:"specifies environmental variables to be loaded into the template"`
	Overwrite bool              `short:"f" long:"force-overwrite" description:"overwrite output file if exists"`
	Verbose   bool              `long:"verbose" description:"verbose mode for debug purposes"`
}

var version string

func main() {
	o := option{
		Values: make(map[string]string),
	}

	if _, err := flags.Parse(&o); err != nil {
		os.Exit(-1)
	}

	if o.Verbose {
		fmt.Printf("ryankurte/utils fcfg version: %s\n", version)
		fmt.Printf("https://github.com/ryankurte/utils/fcfg\n")
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

	if o.Verbose {
		fmt.Printf("Loaded values:\n")
		for k, v := range values {
			fmt.Printf("  - %s:%s\n", k, v)
		}
	}

	if o.Verbose {
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
}
