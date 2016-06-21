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
	json.Unmarshal(body, &m)

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
	ArgsUsage: "<service_id>",
	Description: `Monitors a service while it is being built by its service id.

   Example:
    $ ernest monitor F94034CE-1A57-4A66-AF49-E1E99C5010A2
	`,
	Action: func(c *cli.Context) error {
		_, cfg := setup(c)

		id := c.Args()[0]
		Monitorize(cfg.URL, cfg.Token, id)
		runtime.Goexit()
		return nil
	},
}
