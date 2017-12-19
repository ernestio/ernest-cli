/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helper

import (
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/r3labs/sse"
)

const (
	BUILDCREATE      = "build.create"
	BUILDDELETE      = "build.delete"
	BUILDIMPORT      = "build.import"
	BUILDCREATEDONE  = "build.create.done"
	BUILDDELETEDONE  = "build.delete.done"
	BUILDIMPORTDONE  = "build.import.done"
	BUILDCREATEERROR = "build.create.error"
	BUILDDELETEERROR = "build.delete.error"
	BUILDIMPORTERROR = "build.import.error"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

// Monitorize opens a websocket connection to get input messages
func Monitorize(stream chan *sse.Event) error {
	h := buildhandler{
		writer: uilive.New(),
		stream: stream,
	}

	h.writer.Start()
	defer h.writer.Stop()

	return h.subscribe()
}

// PrintLogs : prints logs inline
func PrintLogs(stream chan *sse.Event) error {
	h := loghandler{stream: stream}
	return h.subscribe()
}

// PrintRawLogs : prints logs inline
func PrintRawLogs(stream chan *sse.Event) error {
	h := rawhandler{stream: stream}
	return h.subscribe()
}
