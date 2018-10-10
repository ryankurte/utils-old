package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mailgun/mailgun-go"
)

type Config struct {
	Domain       string `short:"d" long:"domain" description:"Mailgun domain for sending" env:"MG_DOMAIN"`
	APIKey       string `short:"k" long:"api-key" description:"Mailgun API key" env:"MG_APIKEY"`
	PublicAPIKey string `short:"p" long:"public-api-key" description:"Mailgun public API key" env:"MG_PUBLIC_APIKEY"`

	GetLists  GetLists  `command:"get-lists" description:"Fetch existing mailing lists"`
	AddList   AddList   `command:"create-list" description:"Create a new mailing list"`
	AddMember AddMember `command:"add-member" description:"Add a member to a list"`
	Send      Send      `command:"send" description:"Send an email to a list or address"`

	Version Version `command:"version" description:"Show version and exit"`
	Verbose bool    `short:"v" long:"verbose" description:"Enable verbose logging"`
}

type Send struct {
	Subject string   `short:"s" long:"subject" description:"Email subject" required:"true"`
	Text    string   `short:"b" long:"body-text" description:"Email body text" required:"true"`
	HTML    string   `long:"body-html" description:"Email body in html"`
	From    string   `short:"f" long:"from" description:"Email from address" required:"true"`
	To      []string `short:"t" long:"to" description:"Email to address(es)"`

	CC  []string `long:"cc" description:"Addresses to copy to"`
	BCC []string `long:"bcc" description:"Addresses to blind copy to"`

	Headers map[string]string `short:"h" long:"headers" description:"Email headers"`
	Test    bool              `long:"test" description:"Enable test mode"`
	At      *time.Time        `long:"at" description:"Schedule sending at the provided time"`
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

type AddMember struct {
	List         string                 `short:"l" long:"list" description:"List name" required:"true"`
	NoMerge      bool                   `long:"no-merge" description:"Disables updating existing members"`
	Name         string                 `short:"n" long:"name" description:"Member name" required:"true"`
	Address      string                 `short:"a" long:"address" description:"Member email address" required:"true"`
	Unsubscribed bool                   `long:"unsubscribed" description:"Create member in unsubscribed state"`
	Vars         map[string]interface{} `long:"var" description:"Member variables as K:V pairs"`
}

type Version struct{}

var version string = "NOT SET"

func main() {
	c := Config{}
	p := flags.NewParser(&c, flags.Default)
	_, err := p.Parse()
	if err != nil {
		//fmt.Printf("%s", err)
		os.Exit(-1)
	}

	if c.Verbose {
		fmt.Printf("Connecting to domain: '%s' with public API key: '%s'\n", c.Domain, c.APIKey)
	}

	mg := mailgun.NewMailgun(c.Domain, c.APIKey, c.PublicAPIKey)

	switch p.Active.Name {
	case "get-lists":
		_, lists, err := mg.GetLists(c.GetLists.Limit, c.GetLists.Offset, c.GetLists.Filter)
		if err != nil {
			fmt.Printf("Error fetching list: '%s'\n", err)
			os.Exit(-2)
		}
		fmt.Printf("Lists: %+v", lists)
	case "create-list":
		list := mailgun.List{Address: c.AddList.Address, Name: c.AddList.Name, Description: c.AddList.Description}
		list, err := mg.CreateList(list)
		if err != nil {
			fmt.Printf("Error creating list: '%s'\n", err)
			os.Exit(-3)
		}
		fmt.Printf("Created list")
	case "add-member":
		subscribed, merge := !c.AddMember.Unsubscribed, !c.AddMember.NoMerge
		member := mailgun.Member{Address: c.AddMember.Address, Name: c.AddMember.Name, Subscribed: &subscribed, Vars: c.AddMember.Vars}
		err := mg.CreateMember(merge, c.AddMember.List, member)
		if err != nil {
			fmt.Printf("Error creating member: '%s'\n", err)
			os.Exit(-3)
		}
		fmt.Printf("Created member")
	case "send":
		m := mailgun.NewMessage(c.Send.From, c.Send.Subject, c.Send.Text, c.Send.To...)
		if c.Send.HTML != "" {
			m.SetHtml(c.Send.HTML)
		}
		if c.Send.At != nil {
			m.SetDeliveryTime(*c.Send.At)
		}
		if c.Send.Test {
			m.EnableTestMode()
		}
		for _, v := range c.Send.CC {
			m.AddCC(v)
		}
		for _, v := range c.Send.BCC {
			m.AddBCC(v)
		}
		for k, v := range c.Send.Headers {
			m.AddHeader(k, v)
		}

		status, id, err := mg.Send(m)
		if err != nil {
			fmt.Printf("Error sending message: '%s'\n", err)
			os.Exit(-4)
		}
		fmt.Printf("Message status: '%s' id: '%s'\n", status, id)
	case "version":
		fmt.Printf("%s\n", version)
	default:
		fmt.Printf("Unsupported command\n")
		os.Exit(-5)
	}

}
