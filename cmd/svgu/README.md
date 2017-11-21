# SVGU

A smol utility for programatically generating SVG files using a YML based specification with variable substitution and evaluation.

## Usage

1. Metadata is specified at the top level
2. Variables are specified in the `vars` section then can be referred to by prepending a `$` symbol to the variable name.
3. Equations for evaluation are wrapped by `$(` and `)` and can reference declared variables
4. Entities are declared using these variables and equations to generate SVG entities


### Command line options

```
Usage:
  svgu [OPTIONS]

Application Options:
  -i, --input=  Render configuration file (default: example.yml)
  -o, --output= Render output (default: example.svg)

Help Options:
  -h, --help    Show this help message
```

### Example File

``` yml
---

title: yaml2svg example
description: Example YAML to SVG configuration
width: 1024
height: 840

vars:
  radius: 100

  style-red: {fill: none, stroke: red}
  style-black: {fill: none, stroke: black}
  style-grey: {fill: none, stroke: grey}

entities:
  - {type: circle,  X: $(width / 2),  Y: $(height / 2), R: $radius, Style: $style-red}
  - {type: line,    x1: 0, y1: 0, x2: $width, y2: $height, style: $style-black}
  - {type: ellipse, x: $(width / 2),  y: $(height / 2), w: $(width / 6), h: $(height / 6), style: $style-red}
  - {type: rect,    x: 10, y: 10, w: $(width - 20), h: $(height - 20), style: {fill: none, stroke: red, stroke-width: 4px}}
  - {type: grid, x: 0, y: 0, w: $width, h: $height, n: 64, style: {fill: none, stroke: grey, opacity: 0.9}}
    

```

