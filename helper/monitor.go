/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/r3labs/sse"
)

const (
	// SERVICECREATE ...
	SERVICECREATE = "service.create"
	// SERVICEDELETE ...
	SERVICEDELETE = "service.delete"
	// SERVICEIMPORT ...
	SERVICEIMPORT = "service.import"
	// SERVICECREATEDONE ...
	SERVICECREATEDONE = "service.create.done"
	// SERVICEDELETEDONE ...
	SERVICEDELETEDONE = "service.delete.done"
	// SERVICEIMPORTDONE ...
	SERVICEIMPORTDONE = "service.import.done"
	// SERVICECREATEERROR ...
	SERVICECREATEERROR = "service.create.error"
	// SERVICEDELETEERROR ...
	SERVICEDELETEERROR = "service.delete.error"
	// SERVICEIMPORTERROR ...
	SERVICEIMPORTERROR = "service.import.error"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

// Monitorize opens a websocket connection to get input messages
func Monitorize(host, endpoint, token, stream string) error {
	h := buildhandler{
		writer: uilive.New(),
		stream: OpenStream(host, endpoint, token, stream),
	}

	h.writer.Start()
	defer h.writer.Stop()

	return h.subscribe()
}

// PrintLogs : prints logs inline
func PrintLogs(host, endpoint, token, stream string, blacklist map[string]string) error {
	h := loghandler{
		stream:    OpenStream(host, endpoint, token, stream),
		blacklist: blacklist,
	}

	return h.subscribe()
}

// PrintRawLogs : prints logs inline
func PrintRawLogs(host, endpoint, token, stream string) error {
	h := rawhandler{
		stream: OpenStream(host, endpoint, token, stream),
	}

	return h.subscribe()
}

// OpenStream : opens an sse stream
func OpenStream(host, endpoint, token, stream string) chan *sse.Event {
	ec := make(chan *sse.Event, 1024)

	client := sse.NewClient(host + endpoint)

	//client.EventID = "0"
	client.EncodingBase64 = true
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	go func() {
		err := client.SubscribeChan(stream, ec)
		if err != nil {
			fmt.Println("error connecting to stream: " + err.Error())
			os.Exit(1)
		}
	}()

	return ec
}
