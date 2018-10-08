package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mailgun/mailgun-go"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Domain       string `short:"d" long:"domain" description:"Mailgun domain for sending" required:"true"`
	APIKey       string `short:"k" long:"api-key" description:"Mailgun API key"  required:"true"`
	PublicAPIKey string `short:"p" long:"public-api-key" description:"Mailgun public API key"  required:"true"`

	GetLists GetLists `command:"get-lists"`
	AddList  AddList  `command:"add-list"`
	Send     Send     `command:"send"`

	Verbose bool `short:"v" long:"verbose" description:"Enable verbose logging"`
	Version bool `long:"version" description:"Show version and exit"`
}

type Send struct {
	Subject string `short:"s" long:"subject" description:"Email subject" required:"true"`
	Body    string `short:"b" long:"body" description:"Email body"`
	From    string `short:"f" long:"from" description:"Email from address"`
}

type GetLists struct{}

type AddList struct {
	Name string `long:"name" description:"List name" required:"true"`
}

var version string

func main() {
	c := Config{}
	_, err := flags.Parse(&c)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if c.Version {
		log.Printf("%s\n", version)
		os.Exit(0)
	}

	mg := mailgun.NewMailgun(c.Domain, c.APIKey, c.PublicAPIKey)

	_, lists, err := mg.GetLists(100, 0, "")
	if err != nil {
		log.Fatalf("Error fetching list: %s", err)
	}

	if c.Verbose {
		log.Infof("Found %d lists", len(lists))
	}

}
