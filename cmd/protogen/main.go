package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/ryankurte/utils/cmd/protogen/lib"
)

func main() {
	o := protogen.Options{}
	p := flags.NewParser(&o, flags.Default)
	p.Parse()

	var err error

	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}
}
