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
	"github.com/hokaccha/go-prettyjson"
	"github.com/r3labs/sse"
)

type print func([]byte)

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, endpoint, token, stream string) {
	sseSubscribe(host, endpoint, token, stream, func(body []byte) {
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
	})
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

		yellow := color.New(color.FgYellow).PrintfFunc()
		yellow("%s level=%s user=%s : %s\n", time.Now().Format("02/01/2006 15:04:05"), m.Level, m.User, m.Subject)
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

		fmt.Println(time.Now().Format("02/01/2006 15:04:05"), "level="+m.Level, "user="+m.User, ":", m.Subject+" ", m.Body)
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
