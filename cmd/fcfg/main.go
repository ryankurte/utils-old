package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/jessevdk/go-flags"
)

type option struct {
	Input     flags.Filename    `short:"i" long:"input" description:"input template file (required)"`
	Output    flags.Filename    `short:"o" long:"output" description:"output file (required)"`
	Values    map[string]string `short:"v" long:"values" description:"specifies key:value pairs to be loaded into the template"`
	Keys      []string          `short:"k" long:"keys" description:"specifies environmental variables to be loaded into the template"`
	Config    flags.Filename    `short:"c" long:"config" description:"YAML formatted key-value pairs to be loaded into the template"`
	Overwrite bool              `short:"f" long:"force-overwrite" description:"overwrite output file if exists"`
	Verbose   bool              `long:"verbose" description:"verbose mode for debug purposes"`
	Version   bool              `long:"version" description:"output version tag and exit"`
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
		os.Exit(-2)
	}

	if o.Verbose {
		fmt.Printf("ryankurte/utils fcfg version: %s\n", version)
		fmt.Printf("https://github.com/ryankurte/utils/fcfg\n")
		fmt.Printf("Loading template file: %s\n", o.Input)
	}

	if o.Input == "" || o.Output == "" {
		fmt.Printf("-i,--input and -o,--output arguments are required\n")
		os.Exit(-1)
	}

	f, err := ioutil.ReadFile(string(o.Input))
	if err != nil {
		fmt.Printf("Error opening input file: %s\n", err)
		os.Exit(-3)
	}

	tmpl, err := template.New("").Parse(string(f))
	if err != nil {
		fmt.Printf("Error parsing template: %s\n", err)
		os.Exit(-4)
	}

	values := make(map[string]string)

	// Load environmental keys
	for _, k := range o.Keys {
		values[k] = os.Getenv(k)
	}

	// Load config (if supplied)
	if o.Config != "" {
		c, err := ioutil.ReadFile(string(o.Config))
		if err != nil {
			fmt.Printf("Error loading config gile %s\n", err)
			os.Exit(-5)
		}
		configValues := make(map[string]string)
		yaml.Unmarshal(c, &configValues)
		for k, v := range configValues {
			values[k] = v
		}
	}

	// Load CLI values
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
		os.Exit(-6)
	}

	w := bufio.NewWriter(wr)
	err = tmpl.Execute(w, values)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		os.Exit(-7)
	}

	w.Flush()
	wr.Close()
}
