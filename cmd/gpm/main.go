package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/ryankurte/utils/cmd/gpm/lib"
)

func main() {
	o := gpm.Options{}
	p := flags.NewParser(&o, flags.Default)
	p.Parse()

	log.Printf("Options: %+v", o)

	if p.Active == nil {
		log.Printf("No command selected")
		return
	}

	gpm := gpm.NewGPM(&o.CommonOptions)

	var err error

	switch p.Active.Name {
	case "init":
		err = gpm.Init(&o.Init)
	case "add":
		err = gpm.Add(&o.Add)
	case "sync":
		err = gpm.Sync(&o.Sync)
	case "update":
		err = gpm.Update(&o.Update)
	case "remove":
		err = gpm.Remove(&o.Remove)
	}

	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(-1)
	}
}
