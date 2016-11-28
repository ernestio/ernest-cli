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

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/r3labs/sse"
)

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, token, stream string) {
	url := host + "/events"
	client := sse.NewClient(url)
	client.EncodingBase64 = true
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

// PrettyPrint unmarshals received messages and print it
func PrettyPrint(body []byte) {
	m := model.Message{}

	// clean msg body of any null characters
	cleanedInput := bytes.Trim(body, "\x00")

	err := json.Unmarshal(cleanedInput, &m)
	if err != nil {
		fmt.Println(err)
	}

	if m.Body == "error" || m.Body == "success" {
		os.Exit(0)
	}
	PrintLine(m)
}

// PrintLine prints each received message on a line with its color
func PrintLine(m model.Message) {
	if m.Level == "ERROR" {
		color.Red(m.Body)
	} else if m.Level == "SUCCESS" {
		color.Green(m.Body)
	} else if m.Level == "INFO" {
		color.Yellow(m.Body)
	} else {
		fmt.Println(m.Body)
	}
}
