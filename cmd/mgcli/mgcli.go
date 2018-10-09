package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mailgun/mailgun-go"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Domain       string `short:"d" long:"domain" description:"Mailgun domain for sending" env:"MG_DOMAIN"`
	APIKey       string `short:"k" long:"api-key" description:"Mailgun API key" env:"MG_APIKEY"`
	PublicAPIKey string `short:"p" long:"public-api-key" description:"Mailgun public API key" env:"MG_PUBLIC_APIKEY"`

	GetLists GetLists `command:"get-lists"`
	AddList  AddList  `command:"add-list"`
	Send     Send     `command:"send"`
	Version  Version  `command:"version" description:"Show version and exit"`

	Verbose bool `short:"v" long:"verbose" description:"Enable verbose logging"`
}

type Send struct {
	Subject string   `short:"s" long:"subject" description:"Email subject" required:"true"`
	Body    string   `short:"b" long:"body" description:"Email body" required:"true"`
	From    string   `short:"f" long:"from" description:"Email from address" required:"true"`
	To      []string `short:"t" long:"to" description:"Email to address(es)"`

	Headers map[string]string `short:"h" long:"headers" description:"Email headers"`
}

type limitAndSkip struct {
	Limit  int `long:"limit" description:"Maximum lists to return" default:"20"`
	Offset int `long:"offset" description:"List index offset" default:"0"`
}

type GetLists struct {
	limitAndSkip
	Filter string `long:"filter" description:"List filter"`
}

type AddList struct {
	Address     string `long:"address" description:"List address" required:"true"`
	Name        string `long:"name" description:"List name" required:"true"`
	Description string `long:"description" description:"List description"`
}

type Version struct{}

var version string = "NOT SET"

func main() {
	c := Config{}
	p := flags.NewParser(&c, flags.Default)
	_, err := p.Parse()
	if err != nil {
		log.Fatalf("Invalid arguments")
	}

	mg := mailgun.NewMailgun(c.Domain, c.APIKey, c.PublicAPIKey)

	switch p.Active.Name {
	case "get-lists":
		_, lists, err := mg.GetLists(c.GetLists.Limit, c.GetLists.Offset, c.GetLists.Filter)
		if err != nil {
			log.Fatalf("Error fetching list: %s", err)
		}
		fmt.Printf("Lists: %+v", lists)
	case "add-lists":
		list := mailgun.List{Address: c.AddList.Address, Name: c.AddList.Name, Description: c.AddList.Description}
		list, err := mg.CreateList(list)
		if err != nil {
			log.Fatalf("Error fetching list: %s", err)
		}
		fmt.Printf("Created list")
	case "send":
		m := mailgun.NewMessage(c.Send.From, c.Send.Subject, c.Send.Body, c.Send.To...)
		status, id, err := mg.Send(m)
		if err != nil {
			log.Fatalf("Error sending message: %s", err)
		}
		fmt.Printf("Message status: '%s' id: '%s'", status, id)
	case "version":
		fmt.Printf("%s\n", version)
		os.Exit(0)
	default:
		log.Fatalf("Unsupported command")
	}

}
