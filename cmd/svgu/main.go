package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/ajstarks/svgo/float"
	"github.com/go-yaml/yaml"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/mapstructure"
	"github.com/ryankurte/go-structparse"
)

// Meta is a map of SVG metadata
type Meta map[string]string

// Vars is a map of vars for use in the input file
type Vars map[string]interface{}

// ParseString finds strings that should be evaluated or replaced with objects and replaces
// the instance in the map with an appropriate type/value.
func (v *Vars) ParseString(in string) interface{} {
	if strings.HasPrefix(in, "$(") && strings.HasSuffix(in, ")") {
		// Evaluate strings wrapped in $()
		exp := strings.TrimSuffix(strings.TrimPrefix(in, "$("), ")")

		// Create a new expression and evaluate it
		expression, err := govaluate.NewEvaluableExpression(exp)
		if err != nil {
			fmt.Printf("Expression '%s' error: %s\n", in, err)
			return in
		}
		result, err := expression.Evaluate(*v)
		if err != nil {
			fmt.Printf("Evaluation '%s' error: %s\n", in, err)
			return in
		}

		return result

	} else if strings.HasPrefix(in, "$") {
		// Replace strings starting with $
		key := strings.TrimPrefix(in, "$")
		val, ok := (*v)[key]
		if !ok {
			return in
		}

		return val
	}

	return in
}

// Style is a map of SVG style elements
type Style map[string]string

func (s *Style) String() string {
	str := ""
	for key, val := range *s {
		str += fmt.Sprintf("%s:%s;", key, val)
	}
	return str
}

// Entity is an object to be rendered
type Entity struct {
	Type           string
	X, Y, W, H     float64
	R, N           float64
	X1, Y1, X2, Y2 float64
	Style          Style
}

type options struct {
	Input  flags.Filename `short:"i" long:"input" description:"Render configuration file" default:"example.yml"`
	Output flags.Filename `short:"o" long:"output" description:"Render output" default:"example.svg"`
}

func main() {
	// Load command line options
	o := options{}
	_, err := flags.Parse(&o)
	if err != nil {
		os.Exit(-1)
	}

	// Read the config file into a buffer for further use
	d, err := ioutil.ReadFile(string(o.Input))
	if err != nil {
		os.Exit(-1)
	}

	// Load input file
	// Note that entities are generic at this stage
	c := struct {
		Title       string
		Description string
		Units       string
		Width       float64
		Height      float64
		Vars        Vars
		Meta        Meta
		Entities    []map[string]interface{}
	}{}
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		fmt.Printf("YAML config unmarshalling error: %s\n", err)
		os.Exit(-1)
	}

	// Set width and heigh variables
	c.Vars["width"] = c.Width
	c.Vars["height"] = c.Height

	// Evaluate and replace variables and expressions within generic entities
	structparse.Strings(&c.Vars, &c.Entities)

	// Decode generic entities into useful structures
	entities := make([]Entity, 0)
	mapstructure.Decode(c.Entities, &entities)

	// Create buffer and canvas for rendering
	buff := bytes.NewBuffer(nil)
	canvas := svg.New(buff)

	// Start rendering canvas with pixels and/or units
	if c.Units == "" {
		canvas.Start(c.Width, c.Height)
	} else {
		canvas.Startunit(c.Width, c.Height, c.Units)
	}

	if c.Title != "" {
		canvas.Title(c.Title)
	}

	if c.Description != "" {
		canvas.Desc(c.Description)
	}

	// Render entities
	for _, e := range entities {
		switch e.Type {
		case "circle":
			canvas.Circle(e.X, e.Y, e.R, e.Style.String())
		case "ellipse":
			canvas.Ellipse(e.X, e.Y, e.W, e.H, e.Style.String())
		case "rect":
			canvas.Rect(e.X, e.Y, e.W, e.H, e.Style.String())
		case "center-rect":
			canvas.CenterRect(e.X, e.Y, e.W, e.H, e.Style.String())
		case "grid":
			canvas.Grid(e.X, e.Y, e.W, e.H, e.N, e.Style.String())
		case "line":
			canvas.Line(e.X1, e.Y1, e.X2, e.Y2, e.Style.String())
		}
	}

	canvas.End()

	// Write output file
	err = ioutil.WriteFile(string(o.Output), buff.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %s\n", err)
		os.Exit(-1)
	}

}
