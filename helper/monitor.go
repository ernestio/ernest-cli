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
	"regexp"

	"github.com/ernestio/ernest-cli/model"
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mitchellh/mapstructure"
	"github.com/r3labs/sse"
)

type print func([]byte)

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, endpoint, token, stream string) {
	sseSubscribe(host, endpoint, token, stream, func(event []byte) {
		//		m := model.ServiceNew{}

		// clean msg body of any null characters

		cleanedInput := bytes.Trim(event, "\x00")

		in := make(map[string]interface{})

		err := json.Unmarshal(cleanedInput, &in)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("origin = %#v\n\n", in)

		r := regexp.MustCompile(`^service\.`)
		if r.MatchString(in["_subject"].(string)) {
			m := model.ServiceNew{}

			err := mapstructure.Decode(in, &m)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%#v\n", m)

		} else {
			m := model.ComponentNew{}

			err := mapstructure.Decode(in, &m)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%#v\n", m)

		}

		//		if m.Body == "error" || m.Body == "success" {
		//			os.Exit(0)
		//		}
		//		PrintLine(m)
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

// PrintLine prints each received message on a line with its color
func PrintLine(m interface{}) {
	//func PrintLine(m model.ServiceNew) {
	fmt.Printf("m = %+v\n", m)
	//	if m.Level == "ERROR" {
	//		color.Red(m.Body)
	//	} else if m.Level == "SUCCESS" {
	//		color.Green(m.Body)
	//	} else if m.Level == "INFO" {
	//		color.Yellow(m.Body)
	//	} else {
	//		fmt.Println(m.Body)
	//	}
}
