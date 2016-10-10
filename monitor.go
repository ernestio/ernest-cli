/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/r3labs/sse"
	"github.com/urfave/cli"
)

// Message represents an incomming websocket message
type Message struct {
	Body  string `json:"body"`
	Level string `json:"level"`
}

// NullWriter to disable logging
type NullWriter int

// Write sends to nowhere the log messages
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// PrettyPrint unmarshals received messages and print it
func PrettyPrint(body []byte) {
	m := Message{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return
	}

	if m.Body == "error" || m.Body == "success" {
		os.Exit(0)
	}
	PrintLine(m)
}

// PrintLine prints each received message on a line with its color
func PrintLine(m Message) {
	if m.Level == "ERROR" {
		color.Red(m.Body)
	} else if m.Level == "SUCCESS" {
		color.Green(m.Body)
	} else if m.Level == "INFO" {
		color.Yellow(m.Body)
	} else {
		println(m.Body)
	}
}

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, token, stream string) {
	url := host + "/events"
	client := sse.NewClient(url)
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	err := client.Subscribe(stream, func(msg *sse.Event) {
		if msg.Data != nil {
			PrettyPrint(msg.Data)
		}
	})
	if err != nil {
		log.Println("Failed with: " + err.Error())
		os.Exit(1)
	}
}

// MonitorService command
// Monitorizes an service and shows the actions being performed on it
var MonitorService = cli.Command{
	Name:      "monitor",
	Aliases:   []string{"m"},
	Usage:     "Monitor a service.",
	ArgsUsage: "<service_name>",
	Description: `Monitors a service while it is being built by its service name.

   Example:
    $ ernest monitor my_service
	`,
	Action: func(c *cli.Context) error {
		m, cfg := setup(c)
		if cfg.Token == "" {
			color.Red("You're not allowed to perform this action, please log in")
			return nil
		}

		if len(c.Args()) == 0 {
			color.Red("You should specify an existing service name")
			return nil
		}

		name := c.Args()[0]
		service, err := m.ServiceStatus(cfg.Token, name)
		if err != nil {
			color.Red(err.Error())
			return nil
		}
		parts := strings.Split(service.ID, "-")
		if len(parts) == 0 {
			color.Red("Invalid service specified")
			return nil
		}
		if service.Status == "done" {
			color.Yellow("Service has been successfully built")
			color.Yellow("You can check its information running `ernest-cli service info " + name + "`")
			return nil
		}

		Monitorize(cfg.URL, cfg.Token, parts[len(parts)-1])
		runtime.Goexit()
		return nil
	},
}
