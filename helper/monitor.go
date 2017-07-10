/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mitchellh/mapstructure"
	"github.com/r3labs/sse"
)

type print func([]byte)

var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, endpoint, token, stream string, resc chan string) {
	var s model.ServiceEvent
	var c model.ComponentEvent
	var format string
	var args []interface{}
	var resourceErr error

	writer := uilive.New()
	writer.Start()

	sseSubscribe(host, endpoint, token, stream, func(event []byte) {
		// clean msg body of any null characters
		cleanedInput := bytes.Trim(event, "\x00")

		msg := make(map[string]interface{})

		err := json.Unmarshal(cleanedInput, &msg)
		if err != nil {
			fmt.Println(err)
		}

		subject := msg["_subject"].(string)

		switch subject {
		case "service.create", "service.delete", "service.import":
			s = processServiceEvent(msg)
			format, args = renderOutput(s)
		case "service.create.done", "service.create.error", "service.delete.done", "service.delete.error", "service.import.done", "service.import.error":
			s = processServiceEvent(msg)
			err = renderUpdate(s, model.ComponentEvent{}, args)
			if err != nil {
				log.Println(err)
			}
		default:
			c = processComponentEvent(msg)
			err = renderUpdate(model.ServiceEvent{}, c, args)
			if err != nil {
				resourceErr = err
			}
		}

		time.Sleep(time.Millisecond * 100)
		fmt.Fprintf(writer, format, args...)

		switch subject {
		case "service.create.done":
			writer.Stop()
			resc <- s.Name
		case "service.delete.done", "service.import.done":
			writer.Stop()
			os.Exit(0)
		case "service.create.error", "service.delete.error", "service.import.error":
			writer.Stop()
			fmt.Printf("Message: %s\n\n", red(resourceErr))
			os.Exit(1)
		}
	})
}

func renderUpdate(s model.ServiceEvent, c model.ComponentEvent, a []interface{}) error {
	// component status
	for i, v := range a {
		t := formatType(c.Type)
		if v == t || v == t+"s" {
			if c.Action == "create" && c.State == "running" {
				a[i+3] = yellow("Creating")
			} else if c.Action == "create" && c.State == "completed" {
				a[i+1] = a[i+1].(int) + 1
				if a[i+1] == a[i+2] {
					a[i+3] = green("Created")
				}
			}
			if c.Action == "find" && c.State == "running" {
				a[i+3] = yellow("Searching")
			} else if c.Action == "find" && c.State == "completed" {
				if len(c.Components) > 0 {
					a[i+1] = len(c.Components)
					a[i+2] = len(c.Components)
					a[i+3] = green("Found")
				} else {
					a[i+3] = yellow("None")
				}
			}
			if c.Action == "delete" && c.State == "running" {
				a[i+3] = yellow("Deleting")
			} else if c.Action == "delete" && c.State == "completed" {
				a[i+1] = a[i+1].(int) + 1
				if a[i+1] == a[i+2] {
					a[i+3] = green("Deleted")
				}
			}
			if c.State == "errored" {
				a[i+3] = red("Error")
				return errors.New(c.Error)
			}
		}
	}

	// overall status
	switch s.Subject {
	case "service.delete.done":
		a[len(a)-1] = green("Destroyed")
	case "service.create.done":
		a[len(a)-1] = green("Applied")
	case "service.import.done":
		a[len(a)-1] = green("Imported")
	case "service.create.error", "service.delete.error", "service.import.error":
		a[len(a)-1] = red("Error")
	}

	return nil
}

func renderOutput(s model.ServiceEvent) (string, []interface{}) {
	var blue = color.New(color.FgBlue).SprintFunc()

	f := "\nService Name: %s\n"
	a := []interface{}{blue(s.Name)}
	f = f + "Service ID: %s\n\n"
	a = append(a, blue(s.ID))

	if len(s.Changes) == 0 {
		f = f + green("No changes detected\n")
	} else {
		changes := ParseChanges(s.Changes)

		keys := []string{}
		for key := range changes {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		longestKey := 0
		for _, k := range keys {
			if len(k) > longestKey {
				longestKey = len(k)
			}
		}

		for _, k := range keys {
			f = f + "%-" + strconv.Itoa(longestKey+1) + "s %3d/%-3d %s\n"
			t := formatType(k)
			if s.Subject == "service.import" {
				a = append(a, t, 0, 0, "")
			} else {
				a = append(a, t+"s", 0, changes[k], "")
			}
		}
	}

	f = f + "\nStatus: %s\n\n"
	var status string
	switch s.Subject {
	case "service.create":
		status = yellow("Applying")
	case "service.delete":
		status = yellow("Destroying")
	case "service.import":
		status = yellow("Importing")
	default:
		status = yellow("Unknown")
	}
	a = append(a, status)

	return f, a
}

func ParseChanges(c []model.ComponentEvent) map[string]int {
	seen := map[string]int{}
	for _, v := range c {
		seen[v.Type] += 1
	}
	return seen
}

func formatType(t string) string {
	t = strings.Replace(t, "_", " ", -1)
	return strings.Title(t)
}

func processServiceEvent(s map[string]interface{}) model.ServiceEvent {
	m := model.ServiceEvent{}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &m,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(s)
	if err != nil {
		panic(err)
	}

	return m
}

func processComponentEvent(c map[string]interface{}) model.ComponentEvent {
	m := model.ComponentEvent{}
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &m,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(c)
	if err != nil {
		panic(err)
	}

	return m
}

// PrintLogs : prints logs inline
func PrintLogs(host, endpoint, token, stream string) {
	sseSubscribe(host, endpoint, token, stream, func(body []byte) {
		m := model.Message{}

		// clean msg body of any null characters
		cleanedInput := bytes.Trim(body, "\x00")

		if err := json.Unmarshal(cleanedInput, &m); err != nil {
			fmt.Println(err)
		}

		color.Yellow(m.Subject)
		if len(m.Body) > 0 {
			message, _ := prettyjson.Format([]byte(m.Body))
			fmt.Println(string(message))
		} else {
			fmt.Println("-- Empty string --")
		}
	})
}

// PrintRawLogs : prints logs inline
func PrintRawLogs(host, endpoint, token, stream string) {
	sseSubscribe(host, endpoint, token, stream, func(body []byte) {
		m := model.Message{}

		// clean msg body of any null characters
		cleanedInput := bytes.Trim(body, "\x00")

		if err := json.Unmarshal(cleanedInput, &m); err != nil {
			fmt.Println(err)
		}

		fmt.Println("[" + m.Subject + "] : " + m.Body)
	})
}

func sseSubscribe(host, endpoint, token, stream string, fn print) {
	url := host + endpoint
	client := sse.NewClient(url)
	client.EncodingBase64 = true
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	err := client.Subscribe(stream, func(msg *sse.Event) {
		if msg.Data != nil {
			fn(msg.Data)
		}
	})
	if err != nil {
		log.Println("Failed with: " + err.Error())
		os.Exit(1)
	}
}
