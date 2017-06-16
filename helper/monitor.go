/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mitchellh/mapstructure"
	"github.com/r3labs/sse"
)

type print func([]byte)

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, endpoint, token, stream string) {
	var s model.ServiceEvent
	var c model.ComponentEvent
	var f string
	var a []interface{}

	go sseSubscribe(host, endpoint, token, stream, func(event []byte) {
		// clean msg body of any null characters
		cleanedInput := bytes.Trim(event, "\x00")

		in := make(map[string]interface{})

		err := json.Unmarshal(cleanedInput, &in)
		if err != nil {
			fmt.Println(err)
		}

		subject := in["_subject"].(string)

		switch {
		case subject == "service.create":
			s = processServiceEvent(in)
			f, a = renderOutput(s)
		case subject == "service.create.done":
			s = processServiceEvent(in)
			a = renderUpdate(s, c, a)
		default:
			c = processComponentEvent(in)
			a = renderUpdate(s, c, a)
		}
	})

	writer := uilive.New()
	writer.Start()

	for {
		// Revisit this
		time.Sleep(time.Second * 1)
		fmt.Fprintf(writer, f, a...)

		if s.Subject == "service.create.done" {
			writer.Stop()
			break
		}

		// why doesn't this work?
		//		switch s.Subject {
		//		case "service.create.done", "service.delete.done", "service.import.done", "service.create.error", "service.delete.error", "service.import.error":
		//			writer.Stop()
		//			break
		//		}
	}
	os.Exit(0)
}

func renderUpdate(s model.ServiceEvent, c model.ComponentEvent, a []interface{}) []interface{} {
	var green = color.New(color.FgGreen).SprintFunc()
	var yellow = color.New(color.FgYellow).SprintFunc()
	var red = color.New(color.FgRed).SprintFunc()

	if len(s.Changes) > 0 {
		// component status
		for i, v := range a {
			if v == c.Type {
				switch c.State {
				case "completed":
					a[i+1] = green(c.State)
				case "errored":
					a[i+1] = red(c.State)
					//			errMsg = c.Error
				default:
					a[i+1] = yellow(c.State)
				}
			}
		}
	}

	// overall status
	switch s.Subject {
	case "service.create.done", "service.delete.done", "service.import.done":
		a[len(a)-1] = green("Applied")
	case "service.create.error", "service.delete.error", "service.import.error":
		a[len(a)-1] = red("Error")
		//		f = f + "Message: %s\n\n"
		//		a = append(a, red(errMsg))
	default:
		a[len(a)-1] = yellow("Applying")
	}

	return a
}

func renderOutput(s model.ServiceEvent) (string, []interface{}) {
	var blue = color.New(color.FgBlue).SprintFunc()

	f := "Service ID: %s\n\n"
	a := []interface{}{blue(s.ID)}

	if len(s.Changes) > 0 {
		for _, sc := range s.Changes {
			f = f + "%s...  %s\n"
			a = append(a, sc.Type)
			a = append(a, "")
		}
	}

	f = f + "\nStatus: %s\n"
	a = append(a, "")

	return f, a
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
