package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
)

type Options struct {
	Ports []string       `short:"p" long:"ports" description:"Serial port(s) to connect to" default:"/dev/ttyUSB0"`
	Baud  int            `short:"b" long:"baud" description:"Baud rate for connections" default:"115200"`
	Log   flags.Filename `short:"l" long:"log-file" description:"Log file to write serial outputs to"`
}

type InputLine struct {
	timestamp time.Time
	port      string
	data      string
}

func main() {
	// Load and parse options
	o := Options{}
	_, err := flags.Parse(&o)
	if err != nil {
		os.Exit(-1)
	}

	// Clear logger flags
	log.SetFlags(0)

	log.Printf("Opening port(s): %+v", o.Ports)

	// Open serial ports
	connections := make([]*serial.Port, len(o.Ports))
	for i, p := range o.Ports {
		c := &serial.Config{Name: p, Baud: o.Baud, ReadTimeout: time.Second}
		s, err := serial.OpenPort(c)
		if err != nil {
			log.Fatal(err)
		}

		connections[i] = s
	}

	// Create waitgroups and context for goroutine management
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Line buffer for received data
	lines := make(chan InputLine)

	// Create listener goroutines
	for i := range o.Ports {
		go func(ctx context.Context, name string, port *serial.Port) {
			wg.Add(1)
			defer wg.Done()

			// Buffered reader for line-by-line display
			reader := bufio.NewReader(port)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					b, _, err := reader.ReadLine()
					if err != nil && err != io.EOF {
						log.Printf("Port '%s' error: %s", name, err)
						cancel()
						return
					}
					if len(b) != 0 {
						for _, v := range strings.Split(string(b), "\r\n") {
							lines <- InputLine{time.Now(), name, v}
						}
					}

				}
			}
		}(ctx, o.Ports[i], connections[i])
	}

	// Create interrupt channel
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Main loop prints output and writes to file
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-c:
			break loop
		case line := <-lines:
			log.Printf("[%s %s] %s\n", line.timestamp.Format(time.RFC3339), line.port, line.data)
		}
	}

	cancel()
	wg.Wait()
}
