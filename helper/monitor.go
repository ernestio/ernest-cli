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

	//	writer.Stop()

	//	renderOutput()

	go sseSubscribe(host, endpoint, token, stream, func(event []byte) {

		// clean msg body of any null characters
		cleanedInput := bytes.Trim(event, "\x00")

		in := make(map[string]interface{})

		err := json.Unmarshal(cleanedInput, &in)
		if err != nil {
			fmt.Println(err)
		}

		r := regexp.MustCompile(`^service\.`)
		if r.MatchString(in["_subject"].(string)) {
			s = processServiceEvent(in)
		} else {
			c = processComponentEvent(in)
		}
	})

	writer := uilive.New()

	writer.Start()
	//	defer writer.Stop()

	for {
		time.Sleep(time.Second * 1)

		f, a := parseService(s)

		//		f = f + fs

		//		for _, v := range s.Changes {
		//			f = f + "%s\n"
		//fmt.Printf("this: %#v\n", v.Type)
		//		}

		//		fmt.Printf("f == %#v\n", f)

		//		fmt.Fprintf(writer, f, s.ID, s.Subject, list)
		fmt.Fprintf(writer, f, a...)

		//		fmt.Printf("%#v", s)
		//		if finished {
		//			break
		//		}
		if s.Subject == "service.create.done" {
			writer.Stop()
			break
		}
	}
	//			s.Subject == "service.create.error" ||
	//			s.Subject == "service.delete.done" ||
	//			s.Subject == "service.delete.error" ||
	//			s.Subject == "service.import.done" ||
	//			s.Subject == "service.import.error" {
	//			break
	//		}

	//}
	os.Exit(0)

}

func parseService(s model.ServiceEvent) (string, []interface{}) {

	//	var f string
	f := "ID: %s\n"
	a := []interface{}{s.ID}

	f = f + "Subject: %s\n"
	a = append(a, s.Subject)

	for _, c := range s.Changes {
		f = f + "%s  %s\n"
		a = append(a, c.Type)
		a = append(a, c.State)

		//fmt.Printf("this: %#v\n", v.Type)
	}

	//	vals := make([]interface{}, len(s))
	//	for i, v := range s {
	//		vals[i] = v.Name
	//	}

	//	fmt.Printf("f = %#v\n", f)
	//	fmt.Printf("vals = %#v\n", vals)

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

	//	fmt.Printf("%#v\n", m)
	return m

	// // release CLI
	// if m.Subject == "service.create.done" ||
	// 	m.Subject == "service.create.error" ||
	// 	m.Subject == "service.delete.done" ||
	// 	m.Subject == "service.delete.error" ||
	// 	m.Subject == "service.import.done" ||
	// 	m.Subject == "service.import.error" {
	// 	os.Exit(0)
	// }

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

	//	fmt.Printf("%#v\n", m)
	return m

}

func renderOutput() {

	fmt.Println("RENDER OUTPUT")

	//	writer := uilive.New()

	//	writer.Start()

	//	for i := 0; i <= 60; i++ {
	//		fmt.Fprintf(writer, "Datacenter: %s\n", x)
	//		if finished {
	//			break
	//		}
	//		time.Sleep(time.Second)
	//	}

	//	writer.Stop()

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
