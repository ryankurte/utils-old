package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type option struct {
	Input     flags.Filename    `short:"i" long:"input" description:"input template file (required)"`
	Output    flags.Filename    `short:"o" long:"output" description:"output file (required)"`
	Values    map[string]string `short:"v" long:"values" description:"specifies key:value pairs to be loaded into the template"`
	Keys      []string          `short:"k" long:"keys" description:"specifies environmental variables to be loaded into the template"`
	Configs   []flags.Filename  `short:"c" long:"configs" description:"YAML formatted key-value files to be loaded into the template"`
	Overwrite bool              `short:"f" long:"force-overwrite" description:"overwrite output file if exists"`
	Quiet     bool              `long:"quiet" description:"Quiet mode disables non-error outputs"`
	Version   bool              `long:"version" description:"Output version tag and exit"`
	NoColor   bool              `long:"no-color" description:"Disables colored terminal outputs"`
}

func (o option) Usage() string {
	return "fcfg --input=openvpn.conf.tmpl --output=openvpn.conf --values=ca:ca.crt --values=key:client.key"
}

var version string

func main() {
	log.SetFlags(0)

	// Create option struct and load command line arguments
	o := option{
		Values: make(map[string]string),
	}
	if _, err := flags.Parse(&o); err != nil {
		os.Exit(-1)
	}

	// Disable colour if set
	if o.NoColor {
		color.NoColor = true
	}

	// Print version and exit if --verson is specified
	if o.Version {
		log.Printf("%s\n", version)
		os.Exit(0)
	}

	// Print app name and version if --quiet is not specified
	if !o.Quiet {
		log.Printf("ryankurte/utils fcfg version: %s\n", version)
		log.Printf("https://github.com/ryankurte/utils\n")
	}

	if o.Input == "" || o.Output == "" {
		log.Fatalf(color.RedString("Missing input template file (-i, --input) and/or output file (-o, --output) arguments"))
		os.Exit(-2)
	}

	if !o.Quiet {
		log.Printf(color.CyanString("Loading template file: %s", o.Input))
	}

	if o.Input == "" || o.Output == "" {
		fmt.Printf("-i,--input and -o,--output arguments are required\n")
		os.Exit(-3)
	}

	f, err := ioutil.ReadFile(string(o.Input))
	if err != nil {
		log.Fatalf(color.RedString("Error opening input file: %s", err))
		os.Exit(-4)
	}

	tmpl, err := template.New("").Parse(string(f))
	if err != nil {
		log.Fatalf(color.RedString("Error parsing template: %s", err))
		os.Exit(-5)
	}

	values := make(map[string]string)

	// Load environmental keys
	for _, k := range o.Keys {
		values[k] = os.Getenv(k)
	}

	// Load config files (if supplied)
	for _, v := range o.Configs {
		c, err := ioutil.ReadFile(string(v))
		if err != nil {
			fmt.Printf("Error loading config gile %s\n", err)
			os.Exit(-6)
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

	if !o.Quiet {
		log.Printf(color.CyanString("Loaded values:"))
		for k, v := range values {
			log.Printf(color.BlueString("  - %s:%s", k, v))
		}
	}

	if !o.Quiet {
		log.Printf(color.CyanString("Writing output file: %s", o.Output))
	}

	wr, err := os.Create(string(o.Output))
	if err != nil {
		log.Fatalf(color.RedString("Error creating output file: %s", err))
		os.Exit(-7)
	}

	w := bufio.NewWriter(wr)
	err = tmpl.Execute(w, values)
	if err != nil {
		log.Fatalf(color.RedString("Error executing template: %s", err))
		os.Exit(-8)
	}

	w.Flush()
	wr.Close()

	if !o.Quiet {
		log.Printf(color.CyanString("Output file written"))
	}
}
